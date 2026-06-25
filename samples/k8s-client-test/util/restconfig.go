package util

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func CreateRestConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		log.Printf("using configuration from '%s'", kubeconfig)
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if kubeconfigEnv := os.Getenv("KUBECONFIG"); kubeconfigEnv != "" {
		log.Printf("using configuration from KUBECONFIG '%s'", kubeconfigEnv)
		return clientcmd.BuildConfigFromFlags("", kubeconfigEnv)
	}

	if home := homedir.HomeDir(); home != "" {
		defaultKubeconfig := filepath.Join(home, ".kube", "config")
		if _, err := os.Stat(defaultKubeconfig); err == nil {
			log.Printf("using configuration from default path '%s'", defaultKubeconfig)
			return clientcmd.BuildConfigFromFlags("", defaultKubeconfig)
		} else if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("check default kubeconfig failed: %w", err)
		}
	}

	log.Printf("using in-cluster configuration")
	return rest.InClusterConfig()
}

func CreateRestConfigFromBytes(kubeconfigBytes []byte) (*rest.Config, error) {
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		return nil, err
	}

	return clientConfig.ClientConfig()
}

// CreateMergedRawConfigFromBytes demonstrates the extra capability provided
// by OverridingClientConfig: getting the merged raw kubeconfig.
func CreateMergedRawConfigFromBytes(kubeconfigBytes []byte) (clientcmdapi.Config, error) {
	clientConfig, err := clientcmd.NewClientConfigFromBytes(kubeconfigBytes)
	if err != nil {
		return clientcmdapi.Config{}, err
	}

	return clientConfig.MergedRawConfig()
}
