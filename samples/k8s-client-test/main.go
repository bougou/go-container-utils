package main

import (
	"flag"
	"fmt"

	"github.com/bougou/go-container-utils/samples/k8s-client-test/one"
	"github.com/bougou/go-container-utils/samples/k8s-client-test/util"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "", "(optional) absolute path to the kubeconfig file")
	flag.Parse()

	config, err := util.CreateRestConfig(*kubeconfig)
	if err != nil {
		panic(fmt.Sprintf("build config failed, err: %s", err))
	}

	namespace := "default"
	crdGroup := "openbayes.com"
	crdVersion := "v1alpha1"
	crdResource := "bayesjobs" // Resource name is plural generally.

	if err := testClientset(config, namespace); err != nil {
		panic(err)
	}

	if err := testDynamicClient(config, crdGroup, crdVersion, crdResource, namespace); err != nil {
		panic(err)
	}

	if err := testRESTClient(config, crdGroup, crdVersion, crdResource, namespace); err != nil {
		panic(err)
	}

	one, err := one.NewOne(*kubeconfig)
	if err != nil {
		panic(err)
	}

	nodes, err := one.GetNodes()
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %d nodes\n", len(nodes))

	workers, err := one.GetWorkers()
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %d workers\n", len(workers))

	masters, err := one.GetMasters()
	if err != nil {
		panic(err)
	}
	fmt.Printf("got %d masters\n", len(masters))
}
