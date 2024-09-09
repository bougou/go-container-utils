package main

import (
	"flag"
	"fmt"
	"os"

	ctutils "github.com/bougou/go-container-utils"
	"github.com/kr/pretty"
	"github.com/vishvananda/netlink"
)

func main() {
	var containerID string
	flag.StringVar(&containerID, "cid", "", "container id")
	flag.Parse()

	if containerID == "" {
		fmt.Println("Error, must provide container id")
		os.Exit(1)
	}

	container, err := ctutils.NewContainer(containerID)
	if err != nil {
		panic(fmt.Errorf("load container failed, err: %s", err))
	}

	intfs, links, err := container.GetInterfaces()
	if err != nil {
		panic(fmt.Errorf("get interfaces failed, err: %s", err))
	}

	fmt.Printf("got %d interfaces, %d links", len(intfs), len(links))
	pretty.Println(intfs)
	pretty.Println(links)

	for _, intf := range intfs {
		fmt.Println(intf.Name)
	}

	m, err := container.GetInterfacesNodeMapping()
	if err != nil {
		panic(fmt.Errorf("get interfaces node mapping failed, err: %s", err))
	}
	pretty.Println(m)

	for containerInterface, nodeInterface := range m {
		link, err := netlink.LinkByName(nodeInterface)
		if err != nil {
			panic(fmt.Errorf("get node link (%s) for container (%s) failed, err: %s", nodeInterface, containerInterface, err))
		}
		pretty.Println(link)
	}
}
