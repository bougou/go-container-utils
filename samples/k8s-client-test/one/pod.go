package one

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPods return pods of the specified namespace.
// Empty string means all namespace.
func (one *One) GetPods(namespace string) ([]corev1.Pod, error) {
	// here, use one's clientset capability to access K8s built-in resources.
	client := one.CoreV1().Pods(namespace)

	ctx := context.Background()
	pods, err := client.List(ctx, metav1.ListOptions{})
	if err != nil {
		msg := fmt.Sprintf("Get pod failed, err: %s", err)
		return nil, errors.New(msg)
	}

	return pods.Items, nil
}

func (one *One) GetPodsByFieldSelector(namespace, fieldSelector string) ([]corev1.Pod, error) {
	const limit int64 = 100

	ctx := context.Background()
	client := one.CoreV1().Pods(namespace)

	var allPods []corev1.Pod
	// Continue on the request side means "start from the first page".
	// Pass "" for the initial request; pass the token from the previous response for subsequent pages.
	continueToken := ""

	for {
		pods, err := client.List(ctx, metav1.ListOptions{
			FieldSelector: fieldSelector,
			Limit:         limit,
			Continue:      continueToken,
		})
		if err != nil {
			return nil, err
		}

		allPods = append(allPods, pods.Items...)

		// Continue on the response side is "" when this is the last page (or all items fit in one page).
		// A non-empty value means more items exist; pass it as Continue in the next request.
		continueToken = pods.Continue
		if continueToken == "" {
			break
		}
	}

	return allPods, nil
}
