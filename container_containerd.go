package container

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/defaults"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/errdefs"
	"github.com/vishvananda/netlink"
)

type ContainerdContainer struct {
	ID       string
	hostRoot string
}

var _ Container = (*ContainerdContainer)(nil)

func NewContainerdContainer(containerID string) *ContainerdContainer {
	return &ContainerdContainer{
		ID:       containerID,
		hostRoot: "/",
	}
}

func createContainerdClient() (*containerd.Client, error) {
	host := os.Getenv("CONTAINERD_HOST")
	if host == "" {
		host = "/run/containerd/containerd.sock"
	}
	return containerd.New(host)
}

func ContainerdRootDir() (string, error) {
	cli, err := createContainerdClient()
	if err != nil {
		return "", fmt.Errorf("create containerd client failed, err: %s", err)
	}
	defer cli.Close()

	return defaults.DefaultRootDir, nil
}

func (dc *ContainerdContainer) GetInterfaces() ([]net.Interface, []netlink.Link, error) {
	return nil, nil, ErrNotImplemented
}

func (dc *ContainerdContainer) GetInterfacesNodeMapping() (map[string]string, error) {
	return nil, ErrNotImplemented
}

func (cc *ContainerdContainer) GetOverlayDirs() (lowerDir, upperDir, mergedDir string, err error) {

	cli, err := createContainerdClient()
	if err != nil {
		return "", "", "", fmt.Errorf("create containerd client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	containerService := cli.ContainerService()
	c, err := containerService.Get(ctx, cc.ID)
	if err != nil {
		return "", "", "", fmt.Errorf("get containerd container failed, err: %s", err)
	}

	if c.Snapshotter != "overlayfs" {
		return "", "", "", fmt.Errorf("containerd container snapshotter is not overlayfs")
	}

	snapshotterService := cli.SnapshotService(c.Snapshotter)
	mounts, err := snapshotterService.Mounts(ctx, c.SnapshotKey)
	if err != nil {
		return "", "", "", fmt.Errorf("got snapshotter mounts failed, err: %s", err)
	}

	for _, mount := range mounts {
		for _, mountOption := range mount.Options {
			if strings.HasPrefix(mountOption, "lowerdir=") {
				lowerDir = strings.TrimPrefix(mountOption, "lowerdir=")
			}
			if strings.HasPrefix(mountOption, "upperdir=") {
				upperDir = strings.TrimPrefix(mountOption, "upperdir=")
			}
		}
	}

	mergedDir, err = cc.getRootFS()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get merged dir: %s", err)
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

func (cc *ContainerdContainer) IsExist() (bool, error) {
	cli, err := createContainerdClient()
	if err != nil {
		return false, fmt.Errorf("create containerd client failed, err: %s", err)
	}
	defer cli.Close()

	containerService := cli.ContainerService()
	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")
	if _, err := containerService.Get(ctx, cc.ID); err != nil {
		if errdefs.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("inspect containerd container failed, err: %s", err)
	}

	return true, nil
}

func (cc *ContainerdContainer) IsOverlay() (bool, error) {
	cli, err := createContainerdClient()
	if err != nil {
		return false, fmt.Errorf("create containerd client failed, err: %s", err)
	}
	defer cli.Close()

	containerService := cli.ContainerService()
	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")
	c, err := containerService.Get(ctx, cc.ID)
	if err != nil {
		return false, fmt.Errorf("get containerd container failed, err: %s", err)
	}

	return c.Snapshotter != "overlayfs", nil
}

func (cc *ContainerdContainer) LoadImage(imageTarFilePath string) error {
	return ErrNotImplemented
}

func (cc *ContainerdContainer) Pause() error {
	cli, err := createContainerdClient()
	if err != nil {
		return fmt.Errorf("create containerd client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	container, err := cli.LoadContainer(ctx, cc.ID)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	status, err := task.Status(ctx)
	if err != nil {
		return err
	}

	if status.Status != containerd.Paused {
		return task.Pause(ctx)
	}

	return nil
}

func (cc *ContainerdContainer) Unpause() error {
	cli, err := createContainerdClient()
	if err != nil {
		return fmt.Errorf("create containerd client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	container, err := cli.LoadContainer(ctx, cc.ID)
	if err != nil {
		return err
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		return err
	}

	status, err := task.Status(ctx)
	if err != nil {
		return err
	}

	if status.Status == containerd.Paused {
		return task.Resume(ctx)
	}

	return nil
}

func (dc *ContainerdContainer) WithHostRoot(hostRoot string) {
	dc.hostRoot = hostRoot
}

func (cc *ContainerdContainer) getRootFS() (string, error) {
	cli, err := createContainerdClient()
	if err != nil {
		return "", fmt.Errorf("create containerd client failed, err: %s", err)
	}
	defer cli.Close()

	ctx := namespaces.WithNamespace(context.Background(), "k8s.io")

	c, err := cli.LoadContainer(ctx, cc.ID)
	if err != nil {
		return "", fmt.Errorf("load container failed, err: %s", err)
	}

	task, err := c.Task(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to get container task, err: %s", err)
	}
	defer task.Delete(ctx)

	return fmt.Sprintf("/run/containerd/io.containerd.runtime.v2.task/k8s.io/%s/rootfs", cc.ID), nil
}
