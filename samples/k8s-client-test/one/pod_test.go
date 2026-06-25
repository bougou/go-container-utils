package one

import (
	"fmt"
	"testing"
)

func TestGetPodsByFieldSelector(t *testing.T) {
	one := testCreateOne(t)

	fieldSelector := "spec.nodeName=titan-v1,status.phase!=Failed,status.phase!=Succeeded"

	pods, err := one.GetPodsByFieldSelector("default", fieldSelector)
	if err != nil {
		t.Fatalf("failed to get pods by field selector: %v", err)
	}

	for _, pod := range pods {
		fmt.Println(pod.Name)
	}
}
