package utils

import (
	"fmt"
	"net"

	"github.com/containerd/containerd/pkg/netns"
	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
)

func GetInterfaces(netnsPath string) ([]net.Interface, []netlink.Link, error) {
	var interfaces = []net.Interface{}
	var links = []netlink.Link{}

	netNS := netns.LoadNetNS(netnsPath)
	if err := netNS.Do(func(hostNs ns.NetNS) error {
		intfs, err := net.Interfaces()
		if err != nil {
			return fmt.Errorf("get interfaces failed, err: %s", err)
		}

		for _, intf := range intfs {
			link, err := netlink.LinkByName(intf.Name)
			if err != nil {
				return fmt.Errorf("link name for (%s) failed, err: %s", intf.Name, err)
			}

			links = append(links, link)
			interfaces = append(interfaces, intf)
		}
		return nil
	}); err != nil {
		return nil, nil, fmt.Errorf("failed inside ns, err: %s", err)
	}

	return interfaces, links, nil

}

func GetInterfacesNodeMapping(links []netlink.Link) (map[string]string, error) {
	var ret = map[string]string{}
	for _, link := range links {
		parentIndex := link.Attrs().ParentIndex
		if parentIndex != 0 {
			parentLink, err := netlink.LinkByIndex(link.Attrs().ParentIndex)
			if err != nil {
				return nil, fmt.Errorf("call LinkByIndex failed, err: %s", err)
			}
			ret[link.Attrs().Name] = parentLink.Attrs().Name
		}
	}

	return ret, nil
}
