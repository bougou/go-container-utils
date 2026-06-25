package one

import (
	"errors"
	"fmt"

	"github.com/bougou/go-container-utils/samples/k8s-client-test/util"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

// One is a client for Kubernetes which combines the clientset and dynamic client.
type One struct {
	*kubernetes.Clientset
	dynamic.Interface
}

func NewOne(kubeconfig string) (*One, error) {
	config, err := util.CreateRestConfig(kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		msg := fmt.Sprintf("create kubernetes clientset failed using provided kubeconfig, err: %s", err)
		return nil, errors.New(msg)
	}

	// create the dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		msg := fmt.Sprintf("create kubernetes dynamic client failed using provided kubeconfig, err: %s", err)
		return nil, errors.New(msg)
	}

	one := &One{clientset, dynamicClient}
	return one, nil
}
