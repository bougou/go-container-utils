package one

import (
	"fmt"
	"path/filepath"
	"testing"

	"k8s.io/client-go/util/homedir"
)

func testCreateOne(t *testing.T) *One {
	t.Setenv("KUBECONFIG", filepath.Join(homedir.HomeDir(), ".kube", "config"))

	one, err := NewOne("")
	if err != nil {
		panic(fmt.Sprintf("failed to create one: %v", err))
	}
	return one
}

func TestDescribePod(t *testing.T) {
	one := testCreateOne(t)

	output, err := one.DescribePod("default", "openbayes-server-0")
	if err != nil {
		t.Fatalf("failed to describe pod: %v", err)
	}
	fmt.Println(output)
}
