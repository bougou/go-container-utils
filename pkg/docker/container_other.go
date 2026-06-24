//go:build !linux

package docker

import (
	"net"

	"github.com/bougou/go-container-utils/pkg/errors"
	"github.com/vishvananda/netlink"
)

func (dc *DockerContainer) GetInterfaces() ([]net.Interface, []netlink.Link, error) {
	return nil, nil, errors.ErrNotImplemented
}

func (dc *DockerContainer) GetInterfacesNodeMapping() (map[string]string, error) {
	return nil, errors.ErrNotImplemented
}
