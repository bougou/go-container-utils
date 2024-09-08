package main

import (
	"fmt"

	ctutils "github.com/bougou/go-container-utils"
	"github.com/kr/pretty"
)

func main() {
	containerID := "docker://2e1a4fc87c8fe2c14843f494133773767eea16c14b9060ffc01cddb6292155b6"

	container, err := ctutils.NewContainer(containerID)
	if err != nil {
		panic(err)
	}

	intfs, links, err := container.GetInterfaces()
	if err != nil {
		panic(err)
	}

	fmt.Printf("got %d interfaces, %d links", len(intfs), len(links))
	pretty.Println(intfs)
	pretty.Println(links)

	for _, intf := range intfs {
		fmt.Println(intf.Name)
	}

	m, err := container.GetInterfacesNodeMapping()
	if err != nil {
		panic(err)
	}
	pretty.Println(m)
}
