//go:build !linux

package utils

import (
	"net"

	"github.com/bougou/go-container-utils/pkg/errors"
	"github.com/vishvananda/netlink"
)

func GetInterfaces(netnsPath string) ([]net.Interface, []netlink.Link, error) {
	return nil, nil, errors.ErrNotImplemented
}

func GetInterfacesNodeMapping(links []netlink.Link) (map[string]string, error) {
	return nil, errors.ErrNotImplemented
}
