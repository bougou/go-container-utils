package util

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// HandleUnstructed demonstrates common operations on an unstructured Kubernetes
// object returned by dynamic or REST clients when the API type is not compiled in.
//
// It mutates a deep copy (not the original list item) to show:
//   - metadata accessors: labels, annotations, finalizers
//   - GVK correction via SetGroupVersionKind
//   - spec/status field access via unstructured.Nested* and SetNestedField
//   - partial typed conversion via runtime.DefaultUnstructuredConverter
//
// crdGroup, crdVersion, and crdResource identify the target CRD for logging and
// GVK setup; they are not used for server calls inside this demo.
//
// Alternatives when you know the schema at build time:
//   - Typed client-go clients (e.g. kubernetes.Clientset) or generated CRD clients
//     from controller-gen / kubebuilder — compile-time safety, IDE support
//   - controller-runtime client.Client with registered Scheme types
//   - sigs.k8s.io/yaml or encoding/json into hand-written or generated structs
//
// For unknown or generic CRDs at runtime, unstructured + dynamic.Interface (as
// used by the callers of this function) is the standard client-go approach; the
// Nested* helpers and DefaultUnstructuredConverter are the official apimachinery
// utilities for field access and conversion without generated code.
func HandleUnstructed(obj *unstructured.Unstructured, crdGroup, crdVersion, crdResource string) {
	fmt.Println("\n=== Unstructured Usage Demo (first object) ===")
	fmt.Printf("Target resource: %s.%s/%s\n", crdResource, crdGroup, crdVersion)
	fmt.Printf("Object identity: namespace=%s name=%s uid=%s\n", obj.GetNamespace(), obj.GetName(), obj.GetUID())

	// Use a deep copy for mutation examples so we do not change the original list item.
	working := obj.DeepCopy()

	labels := working.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	labels["demo.k8s.io/handled"] = "true"
	working.SetLabels(labels)

	annotations := working.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}
	annotations["demo.k8s.io/source"] = "handleUnstructed"
	working.SetAnnotations(annotations)

	finalizers := appendIfMissing(working.GetFinalizers(), "demo.k8s.io/finalizer")
	working.SetFinalizers(finalizers)

	working.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   crdGroup,
		Version: crdVersion,
		Kind:    working.GetKind(),
	})
	gvk := working.GroupVersionKind()
	fmt.Printf("GVK after SetGroupVersionKind: %s/%s, Kind=%s\n", gvk.Group, gvk.Version, gvk.Kind)

	phase, found, err := unstructured.NestedString(working.Object, "status", "phase")
	if err != nil {
		fmt.Printf("Read status.phase (string): <error: %v>\n", err)
	} else if found {
		fmt.Printf("Read status.phase (string): %s\n", phase)
	} else {
		fmt.Println("Read status.phase (string): <not found>")
	}

	statusState, found, err := unstructured.NestedString(working.Object, "status", "state")
	if err != nil {
		fmt.Printf("Read status.state (string): <error: %v>\n", err)
	} else if found {
		fmt.Printf("Read status.state (string): %s\n", statusState)
	} else {
		fmt.Println("Read status.state (string): <not found>")
	}

	sshNodePort, found, err := unstructured.NestedInt64(working.Object, "status", "sshNodePort")
	if err != nil {
		fmt.Printf("Read status.sshNodePort (int64): <error: %v>\n", err)
	} else if found {
		fmt.Printf("Read status.sshNodePort (int64): %d\n", sshNodePort)
	} else {
		fmt.Println("Read status.sshNodePort (int64): <not found>")
	}

	lastTransitionTime, found, err := unstructured.NestedMap(working.Object, "status", "lastTransitionTime")
	if err != nil {
		fmt.Printf("Read status.lastTransitionTime (map): <error: %v>\n", err)
	} else if found {
		fmt.Printf("Read status.lastTransitionTime (map): %v\n", lastTransitionTime)
	} else {
		fmt.Println("Read status.lastTransitionTime (map): <not found>")
	}

	if err := unstructured.SetNestedField(working.Object, "updated-by-unstructured-demo", "spec", "demoNote"); err != nil {
		fmt.Printf("Set spec.demoNote failed: %v\n", err)
	}

	specMap, found, err := unstructured.NestedMap(working.Object, "spec")
	if err != nil {
		fmt.Printf("Read spec map failed: %v\n", err)
	} else if found {
		fmt.Printf("Read spec map key count: %d\n", len(specMap))
	}

	replicas, found, err := unstructured.NestedInt64(working.Object, "spec", "replicas")
	if err != nil {
		fmt.Printf("Read spec.replicas (int64): <error: %v>\n", err)
	} else if found {
		fmt.Printf("Read spec.replicas (int64): %d\n", replicas)
	} else {
		fmt.Println("Read spec.replicas (int64): <not found>, set it to 1")
		if err := unstructured.SetNestedField(working.Object, int64(1), "spec", "replicas"); err != nil {
			fmt.Printf("Set spec.replicas failed: %v\n", err)
		}
		replicasAfterSet, foundAfterSet, errAfterSet := unstructured.NestedInt64(working.Object, "spec", "replicas")
		if errAfterSet != nil {
			fmt.Printf("Read spec.replicas (int64): <error: %v>\n", errAfterSet)
		} else if foundAfterSet {
			fmt.Printf("Read spec.replicas (int64): %d\n", replicasAfterSet)
		}
	}

	type objectView struct {
		Metadata struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"metadata"`
	}
	var view objectView
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(working.Object, &view); err != nil {
		fmt.Printf("Convert unstructured -> typed view failed: %v\n", err)
	} else {
		fmt.Printf("Converted typed view: %s/%s\n", view.Metadata.Namespace, view.Metadata.Name)
	}

	fmt.Printf("Labels now: %d, annotations now: %d, finalizers now: %d\n",
		len(working.GetLabels()), len(working.GetAnnotations()), len(working.GetFinalizers()))
}

func appendIfMissing(items []string, target string) []string {
	for _, item := range items {
		if item == target {
			return items
		}
	}
	return append(items, target)
}
