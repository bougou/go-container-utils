package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"

	"github.com/bougou/go-container-utils/samples/k8s-client-test/util"
)

func testDynamicClient(config *rest.Config, crdGroup, crdVersion, crdResource, namespace string) error {
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("create dynamic client failed, err: %w", err)
	}

	gvr := schema.GroupVersionResource{
		Group:    crdGroup,
		Version:  crdVersion,
		Resource: crdResource,
	}
	dynamicList, err := dynamicClient.Resource(gvr).Namespace(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("dynamic client list %s in %s namespace failed, err: %w", crdResource, namespace, err)
	}

	fmt.Printf("\n=== Dynamic Client for CRD %s.%s/%s (namespace: %s, total: %d) ===\n\n", crdResource, crdGroup, crdVersion, namespace, len(dynamicList.Items))
	for _, obj := range dynamicList.Items {
		fmt.Println(obj.GetName())
	}
	if len(dynamicList.Items) > 0 {
		targetName := dynamicList.Items[0].GetName()
		gotObj, err := dynamicClient.Resource(gvr).Namespace(namespace).Get(context.Background(), targetName, metav1.GetOptions{})
		if err != nil {
			return fmt.Errorf("dynamic client get %s/%s in %s namespace failed, err: %w", crdResource, targetName, namespace, err)
		}
		fmt.Printf("Get by name via dynamic client succeeded: %s\n", gotObj.GetName())
		util.HandleUnstructed(&dynamicList.Items[0], crdGroup, crdVersion, crdResource)
	}
	return nil
}
