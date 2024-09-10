package container

import (
	"fmt"
	"net"
	"strings"

	"github.com/vishvananda/netlink"
)

type Container interface {
	WithHostRoot(hostRoot string)
	IsExist() (bool, error)
	IsOverlay() (bool, error)
	GetOverlayDirs() (lowerDir, upperDir, mergeDir string, err error)
	Pause() error
	Unpause() error
	GetInterfaces() ([]net.Interface, []netlink.Link, error)
	GetInterfacesNodeMapping() (map[string]string, error)
}

// runtimeContainerID has the following format:
//   - docker://xxxxxx
//   - containerd://xxxx
func NewContainer(runtimeContainerID string) (Container, error) {
	var runtime string
	var id string

	if strings.HasPrefix(runtimeContainerID, "docker://") {
		runtime = "docker"
		id = strings.TrimPrefix(runtimeContainerID, "docker://")

	} else if strings.HasPrefix(runtimeContainerID, "containerd://") {
		runtime = "containerd"
		id = strings.TrimPrefix(runtimeContainerID, "containerd://")
	}

	switch runtime {
	case "docker":
		return NewDockerContainer(id), nil

	case "containerd":
		return NewContainerdContainer(id), nil

	default:
		return nil, fmt.Errorf("unknown container runtime: (%s)", runtime)
	}
}

func RuntimeRootDir(runtime string) (string, error) {
	switch runtime {
	case "docker":
		return DockerRootDir()

	case "containerd":
		return ContainerdRootDir()

	default:
		return "", fmt.Errorf("unknown container runtime: (%s)", runtime)
	}
}
