// go:build unix

package container

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/pkg/netns"
	"github.com/containernetworking/plugins/pkg/ns"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
	"github.com/vishvananda/netlink"
)

type DockerContainer struct {
	ID string

	hostRoot string
}

var _ Container = (*DockerContainer)(nil)

func NewDockerContainer(containerID string) *DockerContainer {
	return &DockerContainer{
		ID:       containerID,
		hostRoot: "/",
	}
}

func createDockerClient() (*dockerclient.Client, error) {
	return dockerclient.NewClientWithOpts(
		dockerclient.WithAPIVersionNegotiation(),
		dockerclient.WithHost(dockerclient.DefaultDockerHost),
		dockerclient.FromEnv,
	)
}

func DockerRootDir() (string, error) {
	cli, err := createDockerClient()
	if err != nil {
		return "", fmt.Errorf("create docker client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := context.Background()
	info, err := cli.Info(ctx)
	if err != nil {
		return "", fmt.Errorf("get docker info failed, err: %s", err)
	}

	return info.DockerRootDir, nil
}

func (dc *DockerContainer) WithHostRoot(hostRoot string) {
	dc.hostRoot = hostRoot
}

func (dc *DockerContainer) IsExist() (bool, error) {
	cli, err := createDockerClient()
	if err != nil {
		return false, fmt.Errorf("create docker client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := context.Background()
	if _, err := cli.ContainerInspect(ctx, dc.ID); err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("inspect docker container failed, err: %s", err)
	}

	return true, nil
}

func (dc *DockerContainer) IsOverlay() (bool, error) {
	cli, err := createDockerClient()
	if err != nil {
		return false, fmt.Errorf("create docker client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := context.Background()

	c, err := cli.ContainerInspect(ctx, dc.ID)
	if err != nil {
		return false, fmt.Errorf("inspect docker container failed, err: %s", err)
	}

	return c.GraphDriver.Name == "overlay2", nil
}

func (dc *DockerContainer) GetOverlayDirs() (lowerDir, upperDir, mergedDir string, err error) {
	cli, err := createDockerClient()
	if err != nil {
		return "", "", "", fmt.Errorf("create docker client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := context.Background()

	c, err := cli.ContainerInspect(ctx, dc.ID)
	if err != nil {
		return "", "", "", fmt.Errorf("inspect docker container failed, err: %s", err)
	}

	if c.GraphDriver.Name != "overlay2" {
		return "", "", "", fmt.Errorf("docker graph driver is not overlay")
	}

	lowerDir = c.GraphDriver.Data["LowerDir"]
	upperDir = c.GraphDriver.Data["UpperDir"]
	mergedDir = c.GraphDriver.Data["MergedDir"]

	// Remove init layer, left are all image layers.
	dirs := strings.Split(lowerDir, ":")
	if len(dirs) != 0 {
		if strings.HasSuffix(dirs[0], "-init/diff") {
			dirs = dirs[1:]
			lowerDir = strings.Join(dirs, ":")
		}
	}

	if lowerDir == "" {
		err = fmt.Errorf("lower dir can not be empty")
	}
	if upperDir == "" {
		err = fmt.Errorf("upper dir can not be empty")
	}
	if mergedDir == "" {
		err = fmt.Errorf("merged dir can not be empty")
	}

	return
}

func (dc *DockerContainer) Pause() error {
	cli, err := createDockerClient()
	if err != nil {
		return fmt.Errorf("create docker client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := context.Background()

	containerJSON, err := cli.ContainerInspect(ctx, dc.ID)
	if err != nil {
		return fmt.Errorf("inspect docker container failed, err: %s", err)
	}

	if !containerJSON.State.Paused {
		return cli.ContainerPause(ctx, dc.ID)
	}

	return nil
}

func (dc *DockerContainer) Unpause() error {
	cli, err := createDockerClient()
	if err != nil {
		return fmt.Errorf("create docker client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := context.Background()

	containerJSON, err := cli.ContainerInspect(ctx, dc.ID)
	if err != nil {
		return fmt.Errorf("inspect docker container failed, err: %s", err)
	}

	if containerJSON.State.Paused {
		return cli.ContainerUnpause(ctx, dc.ID)
	}

	return nil
}

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
		// for the pause container itself
		// networkMode == "host", means the container is run in host network mode.
		networkContainerID = dc.ID

	default:
		// container id of the associated network container
		// like: "container:2ce8e0caf28d450170d6cfd43087a4a1d0c17f744202271b6ab7e3949e8b9975"
		networkContainerID, _ = strings.CutPrefix(string(networkMode), "container:")
	}

	newtorkContainer, err := cli.ContainerInspect(ctx, networkContainerID)
	if err != nil {
		return nil, nil, fmt.Errorf("inspect docker network container (%s) failed, err: %s", networkContainerID, err)
	}

	var interfaces = []net.Interface{}
	var links = []netlink.Link{}
	// "SandboxKey": "/var/run/docker/netns/5048a1a60e3b",
	// symbolic link on node: /var/run -> /run
	sandboxKey := newtorkContainer.NetworkSettings.SandboxKey

	if strings.HasPrefix(sandboxKey, "/var/run") {
		sandboxKey, _ = strings.CutPrefix(sandboxKey, "/var")
	}

	netnsPath := filepath.Join(dc.hostRoot, sandboxKey)
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

func (dc *DockerContainer) GetInterfacesNodeMapping() (map[string]string, error) {
	_, links, err := dc.GetInterfaces()
	if err != nil {
		return nil, err
	}

	var ret = map[string]string{}
	for _, link := range links {
		parentIndex := link.Attrs().ParentIndex
		if parentIndex != 0 {
			parentLink, err := netlink.LinkByIndex(link.Attrs().ParentIndex)
			if err != nil {
				return nil, err
			}
			ret[link.Attrs().Name] = parentLink.Attrs().Name
		}
	}

	return ret, nil
}
