package fakeclient

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ListPods lists pods in the given namespace via a typed CoreV1 client.
// Accept kubernetes.Interface so the same code works with a real clientset
// or fake.NewSimpleClientset in unit tests.
func ListPods(ctx context.Context, client kubernetes.Interface, namespace string) ([]corev1.Pod, error) {
	podList, err := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList.Items, nil
}

// GetPod fetches a single pod by namespace and name.
func GetPod(ctx context.Context, client kubernetes.Interface, namespace, name string) (*corev1.Pod, error) {
	return client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
}
