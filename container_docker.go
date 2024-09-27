package container

import (
	"context"
	"fmt"
	"os"
	"strings"

	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/errdefs"
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

func (dc *DockerContainer) LoadImage(imageTarFilePath string) error {
	cli, err := createDockerClient()
	if err != nil {
		return fmt.Errorf("create docker client failed, err: %s", err)
	}
	defer cli.Close()

	imageFile, err := os.Open(imageTarFilePath)
	if err != nil {
		return fmt.Errorf("open image tar file failed, err: %s", err)
	}
	defer imageFile.Close()

	loadResponse, err := cli.ImageLoad(context.Background(), imageFile, true)
	if err != nil {
		return fmt.Errorf("load image failed, err: %s", err)
	}
	defer loadResponse.Body.Close()

	return nil
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

func (dc *DockerContainer) WithHostRoot(hostRoot string) {
	dc.hostRoot = hostRoot
}
