package one

import (
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// GetServiceInstances return serviceinstances of the specified namespace.
// Empty string means all namespace.
func (one *One) GetServiceInstances(namespace string) ([]string, error) {

	gvr := schema.GroupVersionResource{
		Group:    "servicecatalog.k8s.io",
		Version:  "v1beta1",
		Resource: "serviceinstances",
	}

	// here, use one's dynamic client capability to access any CRD resources.
	client := one.Resource(gvr).Namespace(namespace)

	ctx := context.Background()
	serviceinstances, err := client.List(ctx, metav1.ListOptions{})
	if err != nil {
		msg := fmt.Sprintf("Get serviceinstances failed, err: %s", err)
		return nil, errors.New(msg)
	}

	out := []string{}
	for _, si := range serviceinstances.Items {
		p := fmt.Sprintf("%s", si)
		out = append(out, p)
	}

	return out, nil
}
