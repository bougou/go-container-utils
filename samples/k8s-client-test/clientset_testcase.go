package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func testClientset(config *rest.Config, namespace string) error {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create clientset failed, err: %w", err)
	}

	pods, err := clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("list pods in %s namespace failed, err: %w", namespace, err)
	}

	fmt.Printf("\n=== ClientSet for CoreV1 Pods (namespace: %s, total: %d) ===\n\n", namespace, len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}
	return nil
}
