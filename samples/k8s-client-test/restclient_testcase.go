package main

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"

	"github.com/bougou/go-container-utils/samples/k8s-client-test/util"
)

func testRESTClient(config *rest.Config, crdGroup, crdVersion, crdResource, namespace string) error {
	cfg := rest.CopyConfig(config)
	cfg.ContentConfig.GroupVersion = &schema.GroupVersion{Group: crdGroup, Version: crdVersion}
	cfg.APIPath = "/apis"
	cfg.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	cfg.UserAgent = rest.DefaultKubernetesUserAgent()

	restClient, err := rest.UnversionedRESTClientFor(cfg)
	if err != nil {
		return fmt.Errorf("create rest client failed, err: %w", err)
	}

	objList := &unstructured.UnstructuredList{}
	if err := restClient.Get().
		Namespace(namespace).
		Resource(crdResource).
		Do(context.Background()).
		Into(objList); err != nil {
		return fmt.Errorf("list %s in %s namespace failed, err: %w", crdResource, namespace, err)
	}

	fmt.Printf("\n=== Rest Client for CRD %s.%s/%s (namespace: %s, total: %d) ===\n\n", crdResource, crdGroup, crdVersion, namespace, len(objList.Items))
	for _, obj := range objList.Items {
		fmt.Println(obj.GetName())
	}
	if len(objList.Items) > 0 {
		targetName := objList.Items[0].GetName()
		gotObj := &unstructured.Unstructured{}
		if err := restClient.Get().
			Namespace(namespace).
			Resource(crdResource).
			Name(targetName).
			Do(context.Background()).
			Into(gotObj); err != nil {
			return fmt.Errorf("get %s/%s in %s namespace by rest client failed, err: %w", crdResource, targetName, namespace, err)
		}
		fmt.Printf("Get by name via rest client succeeded: %s\n", gotObj.GetName())
		util.HandleUnstructed(&objList.Items[0], crdGroup, crdVersion, crdResource)
	}
	return nil
}
