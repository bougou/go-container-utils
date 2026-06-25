package fakeclient

import (
	"context"
	"fmt"
	"testing"

	"github.com/bougou/go-container-utils/samples/k8s-client-test/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// RunFakeClientsetUsageDemo demonstrates client-go fake.Clientset usage:
//   - seed typed objects before tests
//   - exercise business logic through kubernetes.Interface
//   - inspect client.Actions() to verify API calls
func TestFakeClientsetUsage(t *testing.T) {
	t.Helper()

	seedPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx",
			Namespace: "default",
			Labels:    map[string]string{"app": "nginx"},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "nginx", Image: "nginx:1.25"},
			},
		},
	}

	client := fake.NewSimpleClientset(seedPod)
	ctx := context.Background()

	pods, err := ListPods(ctx, client, "default")
	if err != nil {
		t.Fatalf("list pods failed: %v", err)
	}
	if len(pods) != 1 || pods[0].Name != "nginx" {
		t.Fatalf("unexpected list result: %+v", pods)
	}

	gotPod, err := GetPod(ctx, client, "default", "nginx")
	if err != nil {
		t.Fatalf("get pod failed: %v", err)
	}
	if gotPod.Labels["app"] != "nginx" {
		t.Fatalf("unexpected pod labels: %v", gotPod.Labels)
	}

	newPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "redis",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{Name: "redis", Image: "redis:7"},
			},
		},
	}
	if _, err := client.CoreV1().Pods("default").Create(ctx, newPod, metav1.CreateOptions{}); err != nil {
		t.Fatalf("create pod failed: %v", err)
	}

	pods, err = ListPods(ctx, client, "default")
	if err != nil {
		t.Fatalf("list pods after create failed: %v", err)
	}
	if len(pods) != 2 {
		t.Fatalf("want 2 pods after create, got %d", len(pods))
	}

	if err := util.MatchActionsStrict(client.Actions(), []util.ExpectedAction{
		{Verb: "list", Namespace: "default", Resource: "pods"},
		{Verb: "get", Namespace: "default", Resource: "pods", Name: "nginx"},
		{Verb: "create", Namespace: "default", Resource: "pods", Name: "redis"},
		{Verb: "list", Namespace: "default", Resource: "pods"},
	}); err != nil {
		t.Fatal(err)
	}

	t.Logf("fake clientset demo finished, recorded %d actions", len(client.Actions()))
	fmt.Printf("fake clientset demo: %d pods, %d actions recorded\n", len(pods), len(client.Actions()))
}
