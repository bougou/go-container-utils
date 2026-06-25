package fakeclient

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var mongodbGVR = schema.GroupVersionResource{
	Group:    "mongodbcommunity.mongodb.com",
	Version:  "v1",
	Resource: "mongodbcommunities",
}

func GetMongoDB(ctx context.Context, client dynamic.Interface, namespace string, name string) (*unstructured.Unstructured, error) {
	// GET /apis/mongodbcommunity.mongodb.com/v1/namespaces/{namespace}/mongodbcommunity/
	get, err := client.Resource(mongodbGVR).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return get, nil
}

func ListMongoDB(ctx context.Context, client dynamic.Interface, namespace string) ([]unstructured.Unstructured, error) {
	// GET /apis/mongodbcommunity.mongodb.com/v1/namespaces/{namespace}/mongodbcommunity/
	list, err := client.Resource(mongodbGVR).Namespace(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// all namespaces
	// GET /apis/mongodbcommunity.mongodb.com/v1/mongodbcommunity/
	// _, err = client.Resource(mongodbGVR).List(ctx, metav1.ListOptions{})
	// if err != nil {
	// 	return nil, err
	// }

	return list.Items, nil
}
