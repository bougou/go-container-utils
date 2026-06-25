package one

import (
	"k8s.io/kubectl/pkg/describe"
)

func (one *One) DescribePod(namespace, name string) (string, error) {
	// 创建 describer
	podDescriber := describe.PodDescriber{Interface: one.Clientset}

	// 调用 Describe
	output, err := podDescriber.Describe("default", "openbayes-server-0", describe.DescriberSettings{ShowEvents: true})
	if err != nil {
		panic(err)
	}

	return output, nil
}
