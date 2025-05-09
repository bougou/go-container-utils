//go:build !linux

package container

import (
	"net"

	"github.com/vishvananda/netlink"
)

func GetInterfaces(netnsPath string) ([]net.Interface, []netlink.Link, error) {
	return nil, nil, ErrNotImplemented
}

func GetInterfacesNodeMapping(links []netlink.Link) (map[string]string, error) {
	return nil, ErrNotImplemented
}
