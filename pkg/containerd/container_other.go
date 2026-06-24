//go:build !linux

package containerd

import (
	"net"

	"github.com/bougou/go-container-utils/pkg/errors"
	"github.com/vishvananda/netlink"
)

func (cc *ContainerdContainer) GetInterfaces() ([]net.Interface, []netlink.Link, error) {
	return nil, nil, errors.ErrNotImplemented
}

func (cc *ContainerdContainer) GetInterfacesNodeMapping() (map[string]string, error) {
	return nil, errors.ErrNotImplemented
}
