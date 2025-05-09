//go:build !linux

package container

import (
	"net"

	"github.com/vishvananda/netlink"
)

func (cc *ContainerdContainer) GetInterfaces() ([]net.Interface, []netlink.Link, error) {
	return nil, nil, ErrNotImplemented
}

func (cc *ContainerdContainer) GetInterfacesNodeMapping() (map[string]string, error) {
	return nil, ErrNotImplemented
}
