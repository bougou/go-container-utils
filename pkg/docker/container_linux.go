package docker

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/bougou/go-container-utils/pkg/utils"
	"github.com/vishvananda/netlink"
)

func (dc *DockerContainer) GetInterfaces() ([]net.Interface, []netlink.Link, error) {
	cli, err := createDockerClient()
	if err != nil {
		return nil, nil, fmt.Errorf("create docker client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := context.Background()

	c, err := cli.ContainerInspect(ctx, dc.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("inspect docker container failed, err: %s", err)
	}

	var networkContainerID string

	networkMode := c.HostConfig.NetworkMode

	switch networkMode {
	case "none", "host":
		// networkMode == "none", means the container is the pause container itself.
		// networkMode == "host", means the container is run in host network mode.
		networkContainerID = dc.ID

	default:
		// container id of the associated network container
		// like: "container:2ce8e0caf28d450170d6cfd43087a4a1d0c17f744202271b6ab7e3949e8b9975"
		networkContainerID, _ = strings.CutPrefix(string(networkMode), "container:")
	}

	// Get sandbox key file path from the network container.

	newtorkContainer, err := cli.ContainerInspect(ctx, networkContainerID)
	if err != nil {
		return nil, nil, fmt.Errorf("inspect docker network container (%s) failed, err: %s", networkContainerID, err)
	}

	// "SandboxKey": "/var/run/docker/netns/5048a1a60e3b",
	// symbolic link on node: /var/run -> /run
	sandboxKey := newtorkContainer.NetworkSettings.SandboxKey
	if strings.HasPrefix(sandboxKey, "/var/run") {
		sandboxKey, _ = strings.CutPrefix(sandboxKey, "/var")
	}

	netnsPath := filepath.Join(dc.hostRoot, sandboxKey)

	return utils.GetInterfaces(netnsPath)
}

func (dc *DockerContainer) GetInterfacesNodeMapping() (map[string]string, error) {
	_, links, err := dc.GetInterfaces()
	if err != nil {
		return nil, fmt.Errorf("call GetInterfaces failed, err: %s", err)
	}

	return utils.GetInterfacesNodeMapping(links)
}
