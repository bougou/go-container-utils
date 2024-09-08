package container

import (
	"fmt"
	"strings"
)

type Container interface {
	IsExist() (bool, error)
	IsOverlay() (bool, error)
	GetOverlayDirs() (lowerDir, upperDir, mergeDir string, err error)
	Pause() error
	Unpause() error
}

// runtimeID has the following format:
// docker://xxxxxx
// containerd://xxxx
func NewContainer(runtimeID string) (Container, error) {
	var runtime string
	var id string

	if strings.HasPrefix(runtimeID, "docker://") {
		runtime = "docker"
		id = strings.TrimPrefix(runtimeID, "docker://")

	} else if strings.HasPrefix(runtimeID, "containerd://") {
		runtime = "containerd"
		id = strings.TrimPrefix(runtimeID, "containerd://")
	}

	switch runtime {
	case "docker":
		return NewDockerContainer(id), nil

	case "containerd":
		return NewContainerdContainer(id), nil

	default:
		return nil, fmt.Errorf("unknown container runtime: %s", runtime)
	}
}

func RuntimeRootDir(runtime string) (string, error) {
	switch runtime {
	case "docker":
		return DockerRootDir()

	case "containerd":
		return ContainerdRootDir()

	default:
		return "", fmt.Errorf("unknown container runtime: %s", runtime)
	}
}
