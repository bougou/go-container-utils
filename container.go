package container

import (
	"fmt"
	"net"
	"strings"

	"github.com/vishvananda/netlink"
)

type Runtime string

const (
	RuntimeDocker     Runtime = "docker"
	RuntimeContainerd Runtime = "containerd"
)

var ErrNotImplemented error = fmt.Errorf("not implemented")

type Container interface {
	GetInterfaces() ([]net.Interface, []netlink.Link, error)
	GetInterfacesNodeMapping() (map[string]string, error)
	GetOverlayDirs() (lowerDir, upperDir, mergeDir string, err error)
	IsExist() (bool, error)
	IsOverlay() (bool, error)
	LoadImage(imageTarFilePath string) error
	Pause() error
	Unpause() error
	WithHostRoot(hostRoot string)
}

// runtimeContainerID has the following format:
//   - docker://xxxxxx
//   - containerd://xxxx
func NewContainer(runtimeContainerID string) (Container, error) {
	var runtime Runtime
	var id string

	if strings.HasPrefix(runtimeContainerID, "docker://") {
		runtime = RuntimeDocker
		id = strings.TrimPrefix(runtimeContainerID, "docker://")

	} else if strings.HasPrefix(runtimeContainerID, "containerd://") {
		runtime = "containerd"
		id = strings.TrimPrefix(runtimeContainerID, "containerd://")
	}

	switch runtime {
	case RuntimeDocker:
		return NewDockerContainer(id), nil

	case RuntimeContainerd:
		return NewContainerdContainer(id), nil

	default:
		return nil, fmt.Errorf("unknown container runtime: (%s)", runtime)
	}
}

func RuntimeRootDir(runtime Runtime) (string, error) {
	switch runtime {
	case RuntimeDocker:
		return DockerRootDir()

	case RuntimeContainerd:
		return ContainerdRootDir()

	default:
		return "", fmt.Errorf("unknown container runtime: (%s)", runtime)
	}
}
