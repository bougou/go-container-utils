package container

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/containerd/containerd/namespaces"
	"github.com/vishvananda/netlink"
)

func (cc *ContainerdContainer) GetInterfaces() ([]net.Interface, []netlink.Link, error) {
	cli, err := createContainerdClient()
	if err != nil {
		return nil, nil, fmt.Errorf("create containerd client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	c, err := cli.LoadContainer(ctx, cc.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("load containerd container failed, err: %s", err)
	}

	task, err := c.Task(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("get containerd task failed, err: %s", err)
	}

	pid := task.Pid()
	netnsPath := fmt.Sprintf("/proc/%d/ns/net", pid)
	netnsPath = filepath.Join(cc.hostRoot, netnsPath)

	return GetInterfaces(netnsPath)
}

func (cc *ContainerdContainer) GetInterfacesNodeMapping() (map[string]string, error) {
	_, links, err := cc.GetInterfaces()
	if err != nil {
		return nil, fmt.Errorf("call GetInterfaces failed, err: %s", err)
	}

	return GetInterfacesNodeMapping(links)
}
