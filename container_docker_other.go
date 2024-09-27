//go:build !linux

package container

import (
	"net"

	"github.com/vishvananda/netlink"
)

func (dc *DockerContainer) GetInterfaces() ([]net.Interface, []netlink.Link, error) {
	return nil, nil, ErrNotImplemented
}

func (dc *DockerContainer) GetInterfacesNodeMapping() (map[string]string, error) {
	return nil, ErrNotImplemented
}
