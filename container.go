package container

import (
	"fmt"
	"net"
	"strings"

	"github.com/vishvananda/netlink"
)

// Runtime represents the type of container runtime being used
type Runtime string

const (
	// RuntimeDocker represents the Docker container runtime
	RuntimeDocker Runtime = "docker"
	// RuntimeContainerd represents the Containerd container runtime
	RuntimeContainerd Runtime = "containerd"
)

// ErrNotImplemented is returned when a requested operation is not implemented
var ErrNotImplemented error = fmt.Errorf("not implemented")

// Container defines the interface for container operations.
// It provides methods to interact with and retrieve information about containers
// across different container runtimes.
type Container interface {
	// GetInterfaces retrieves the network interfaces and links of the container.
	// The information is fetched from the host's perspective but within the container's
	// network namespace, effectively providing the same view as if queried from inside the container.
	GetInterfaces() ([]net.Interface, []netlink.Link, error)

	// GetInterfacesNodeMapping returns a mapping between container interface names and their
	// corresponding interface names on the host node.
	// Example mapping:
	// {"eth0": "cali97e0633f831"}
	GetInterfacesNodeMapping() (map[string]string, error)

	// GetOverlayDirs returns the paths of the overlay filesystem directories:
	// - lowerDir: the read-only lower layer
	// - upperDir: the writable upper layer
	// - mergeDir: the merged view of the filesystem
	GetOverlayDirs() (lowerDir, upperDir, mergeDir string, err error)

	// IsExist checks if the container currently exists in the runtime
	IsExist() (bool, error)

	// IsOverlay checks if the container is using an overlay filesystem
	IsOverlay() (bool, error)

	// LoadImage loads a container image from a tar file into the runtime
	LoadImage(imageTarFilePath string) error

	// Pause suspends all processes within the container
	Pause() error

	// Unpause resumes all processes within the container
	Unpause() error

	// WithHostRoot sets the host root directory for container operations
	WithHostRoot(hostRoot string)
}

// NewContainer creates a new Container instance based on the provided runtime container ID.
// The runtimeContainerID should be in one of these formats:
//   - docker://<container-id>
//   - containerd://<container-id>
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

// RuntimeRootDir returns the root directory path for the specified container runtime.
// This is where the runtime stores its container data and metadata.
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
