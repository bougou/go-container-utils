package fakeclient

import (
	"context"
	"testing"
	"time"

	"github.com/bougou/go-container-utils/samples/k8s-client-test/util"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
)

func newMongoDBCommunity(name, namespace string, members int64) *unstructured.Unstructured {
	mdb := &unstructured.Unstructured{}
	mdb.SetUnstructuredContent(map[string]interface{}{
		"apiVersion": "mongodbcommunity.mongodb.com/v1",
		"kind":       "MongoDBCommunity",
		"metadata": map[string]interface{}{
			"name":      name,
			"namespace": namespace,
		},
		"spec": map[string]interface{}{
			"members": members,
		},
	})
	return mdb
}

func newPodWithVolumes(name, config, secret string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Pod",
			"metadata": map[string]interface{}{
				"namespace":         "default",
				"name":              name,
				"creationTimestamp": time.Now().Format(time.RFC3339),
			},
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name": name,
					},
				},
			},
			"volumes": []interface{}{
				map[string]interface{}{
					"name": secret,
					"secret": map[string]interface{}{
						"secretName": secret,
					},
				},
				map[string]interface{}{
					"name": config,
					"configMap": map[string]interface{}{
						"name": config,
					},
				},
			},
		},
	}
}

func TestFakeDynamicClient_GetMongoDB_NotFound(t *testing.T) {
	client := fake.NewSimpleDynamicClient(runtime.NewScheme(), newMongoDBCommunity("mongodb-test", "default", 3))

	_, err := GetMongoDB(context.Background(), client, "default", "missing")
	if !apierrors.IsNotFound(err) {
		t.Fatalf("want NotFound, got %v", err)
	}

	if err := util.MatchActionsStrict(client.Actions(), []util.ExpectedAction{
		{
			Verb:      "get",
			Namespace: "default",
			Name:      "missing",
			Resource:  "mongodbcommunities",
		},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFakeDynamicClient_GetMongoDB(t *testing.T) {
	const name = "mongodb-test"

	client := fake.NewSimpleDynamicClient(runtime.NewScheme(), newMongoDBCommunity(name, "default", 3))

	got, err := GetMongoDB(context.Background(), client, "default", name)
	if err != nil {
		t.Fatal(err)
	}
	if got.GetName() != name {
		t.Fatalf("unexpected object name: got %q, want %q", got.GetName(), name)
	}

	if err := util.MatchActionsStrict(client.Actions(), []util.ExpectedAction{
		{
			Verb:      "get",
			Namespace: "default",
			Name:      name,
			Resource:  "mongodbcommunities",
		},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFakeDynamicClient_ListMongoDB(t *testing.T) {
	client := fake.NewSimpleDynamicClient(runtime.NewScheme(), newMongoDBCommunity("mongodb-test", "default", 3))

	items, err := ListMongoDB(context.Background(), client, "default")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 {
		t.Fatalf("want 1 item, got %d", len(items))
	}

	if err := util.MatchActionsStrict(client.Actions(), []util.ExpectedAction{
		{
			Verb:      "list",
			Namespace: "default",
			Resource:  "mongodbcommunities",
		},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestFakeDynamicClient_SeedPodWithVolumes(t *testing.T) {
	client := fake.NewSimpleDynamicClient(
		runtime.NewScheme(),
		newPodWithVolumes("config-pod", "properties", "tokens"),
	)

	podGVR := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	pod, err := client.Resource(podGVR).Namespace("default").Get(
		context.Background(),
		"config-pod",
		metav1.GetOptions{},
	)
	if err != nil {
		t.Fatal(err)
	}
	if pod.GetName() != "config-pod" {
		t.Fatalf("unexpected pod name: got %q, want %q", pod.GetName(), "config-pod")
	}
}
