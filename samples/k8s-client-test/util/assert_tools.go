package util

import (
	"bytes"
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kubetesting "k8s.io/client-go/testing"
)

type ExpectedAction struct {
	Verb      string
	Name      string
	Namespace string
	Resource  string

	// Patch action
	PatchType    types.PatchType
	PatchPayload []map[string]interface{}
}

// MatchActionsStrict verifies got actions 1-to-1 match expected actions in order.
// It returns nil when all actions match, or an error describing the mismatch.
func MatchActionsStrict(got []kubetesting.Action, expected []ExpectedAction) error {
	if len(expected) != len(got) {
		return fmt.Errorf("executed actions number not equal, expected %d, got %d", len(expected), len(got))
	}

	for i, expectedAction := range expected {
		if err := MatchActionStrict(got[i], expectedAction); err != nil {
			return fmt.Errorf("action %d does not match: %w", i, err)
		}
	}

	return nil
}

// MatchActions verifies each expected action can be found in got actions.
// It does not require an exact 1-to-1 ordered match.
func MatchActions(got []kubetesting.Action, expected []ExpectedAction) error {
	if len(expected) > len(got) {
		return fmt.Errorf("executed actions number too short, expected %d, got %d", len(expected), len(got))
	}

	for i, expectedAction := range expected {
		if err := matchExpectedAction(got, expectedAction); err != nil {
			return fmt.Errorf("action %d does not match any of the got actions: %w", i, err)
		}
	}

	return nil
}

// MatchActionStrict checks whether a single got action matches the expected action.
func MatchActionStrict(got kubetesting.Action, expected ExpectedAction) error {
	switch expected.Verb {
	case "get":
		getAction, ok := got.(kubetesting.GetAction)
		if !ok {
			return fmt.Errorf("expected verb is get, but got %T", got)
		}

		if getAction.GetName() != expected.Name {
			return fmt.Errorf("name not matched for GetAction: got %q, want %q", getAction.GetName(), expected.Name)
		}

		return validateNamespaceAndResource(getAction, expected)

	case "list":
		listAction, ok := got.(kubetesting.ListAction)
		if !ok {
			return fmt.Errorf("expected verb is list, but got %T", got)
		}

		return validateNamespaceAndResource(listAction, expected)

	case "watch":
		watchAction, ok := got.(kubetesting.WatchAction)
		if !ok {
			return fmt.Errorf("expected verb is watch, but got %T", got)
		}

		return validateNamespaceAndResource(watchAction, expected)

	case "create":
		createAction, ok := got.(kubetesting.CreateAction)
		if !ok {
			return fmt.Errorf("expected verb is create, but got %T", got)
		}

		if err := validateNamespaceAndResource(createAction, expected); err != nil {
			return err
		}

		if expected.Name == "" {
			return nil
		}

		metaObj, ok := createAction.GetObject().(metav1.Object)
		if !ok {
			return fmt.Errorf("create object does not implement metav1.Object: %T", createAction.GetObject())
		}
		if metaObj.GetName() != expected.Name {
			return fmt.Errorf("name not matched for CreateAction: got %q, want %q", metaObj.GetName(), expected.Name)
		}

		return nil

	case "update":
		updateAction, ok := got.(kubetesting.UpdateAction)
		if !ok {
			return fmt.Errorf("verb is update, but got %T", got)
		}

		return validateNamespaceAndResource(updateAction, expected)

	case "delete":
		deleteAction, ok := got.(kubetesting.DeleteAction)
		if !ok {
			return fmt.Errorf("verb is delete, but got %T", got)
		}

		if deleteAction.GetName() != expected.Name {
			return fmt.Errorf("name not matched for DeleteAction: got %q, want %q", deleteAction.GetName(), expected.Name)
		}

		return validateNamespaceAndResource(deleteAction, expected)

	case "patch":
		patchAction, ok := got.(kubetesting.PatchAction)
		if !ok {
			return fmt.Errorf("verb is patch, but got %T", got)
		}

		if patchAction.GetName() != expected.Name {
			return fmt.Errorf("name not matched for PatchAction: got %q, want %q", patchAction.GetName(), expected.Name)
		}

		if err := validateNamespaceAndResource(patchAction, expected); err != nil {
			return err
		}

		if patchAction.GetPatchType() != expected.PatchType {
			return fmt.Errorf("patch type not matched for PatchAction: got %q, want %q", patchAction.GetPatchType(), expected.PatchType)
		}

		patchBytes, err := json.Marshal(expected.PatchPayload)
		if err != nil {
			return fmt.Errorf("marshal expected patch payload: %w", err)
		}

		if !bytes.Equal(patchAction.GetPatch(), patchBytes) {
			return fmt.Errorf("patch payload not equal for PatchAction: got %s, want %s", patchAction.GetPatch(), patchBytes)
		}

		return nil
	}

	return fmt.Errorf("unexpected verb %q", expected.Verb)
}

func matchExpectedAction(got []kubetesting.Action, expectedAction ExpectedAction) error {
	for _, gotAction := range got {
		if err := MatchActionStrict(gotAction, expectedAction); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no matching action found for verb %q", expectedAction.Verb)
}

func validateNamespaceAndResource(action kubetesting.Action, expectedAction ExpectedAction) error {
	if action.GetNamespace() != expectedAction.Namespace {
		return fmt.Errorf("namespace not matched: got %q, want %q", action.GetNamespace(), expectedAction.Namespace)
	}
	if action.GetResource().Resource != expectedAction.Resource {
		return fmt.Errorf("resource not matched: got %q, want %q", action.GetResource().Resource, expectedAction.Resource)
	}
	return nil
}
