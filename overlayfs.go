package container

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/containerd/containerd/mount"
	"github.com/kr/pretty"
)

// OverlayFS represents a overlay filesystem composed from the given
// lower dir and upper dir and merged dir.
type OverlayFS struct {
	lowerDir  string
	upperDir  string
	mergedDir string

	m *mount.Mount
}

func NewOverlayFS(lowerDir, upperDir, mergedDir string) *OverlayFS {
	return &OverlayFS{
		lowerDir:  lowerDir,
		upperDir:  upperDir,
		mergedDir: mergedDir,
	}
}

func (fs *OverlayFS) workDir() string {
	parent := filepath.Dir(fs.upperDir)
	return path.Join(parent, "work")
}

// mount_ro constructs and returns a readonly mount.Mount
func (fs *OverlayFS) mount_ro() mount.Mount {
	lowerDirs := strings.Split(fs.lowerDir, ":")
	if fs.upperDir != "" {
		lowerDirs = append([]string{fs.upperDir}, lowerDirs...)
	}

	mountOptions := []string{
		"ro",
		"relatime",
		fmt.Sprintf("lowerdir=%s", strings.Join(lowerDirs, ":")),
	}

	return mount.Mount{
		Type:    "overlay",
		Source:  "overlay",
		Target:  fs.mergedDir,
		Options: mountOptions,
	}
}

// mount constructs and returns a read/write mount.Mount
func (fs *OverlayFS) mount() mount.Mount {
	mountOptions := []string{
		"rw",
		"relatime",
		fmt.Sprintf("lowerdir=%s", fs.lowerDir),
		fmt.Sprintf("upperdir=%s", fs.upperDir),
		fmt.Sprintf("workdir=%s", fs.workDir()),
	}

	return mount.Mount{
		Type:    "overlay",
		Source:  "overlay",
		Target:  fs.mergedDir,
		Options: mountOptions,
	}
}

func (fs *OverlayFS) Mount(readOnly bool) error {
	if fs.lowerDir == "" {
		return fmt.Errorf("overlayfs lower dir can not be empty")
	}

	var m mount.Mount

	if readOnly {
		if fs.upperDir == "" {
			return fmt.Errorf("overlayfs upper dir can not be empty")
		}

		// the workDir directory will be automatically created for read/write mount.Mount is mounted.
		// So, mkdir is only needed for readonly Mount.
		if err := os.MkdirAll(fs.workDir(), 0711); err != nil {
			return fmt.Errorf("failed to create overlayfs work dir (%s): %s", fs.workDir(), err)
		}

		m = fs.mount_ro()
	} else {
		m = fs.mount()
	}

	// store the reference to mount.Mount in order to umount
	fs.m = &m

	if err := os.MkdirAll(fs.mergedDir, 0711); err != nil {
		return fmt.Errorf("failed to create overlayfs merged dir (%s): %s", fs.mergedDir, err)
	}

	if err := m.Mount(""); err != nil {
		pretty.Println(m)
		return err
	}

	return nil
}

func (fs *OverlayFS) Unmount() error {
	if fs.m == nil {
		return nil
	}

	m := *fs.m
	if err := mount.UnmountMounts([]mount.Mount{m}, "", 0); err != nil {
		pretty.Println(m)
		return err
	}

	return nil
}

// Clear removes the work dir and merged dir. It does not removes the upper dir and lower dir.
func (fs *OverlayFS) Clear() error {
	if err := os.RemoveAll(fs.mergedDir); err != nil {
		return fmt.Errorf("failed to remove overlayfs merged dir (%s): %s", fs.mergedDir, err)
	}

	if err := os.RemoveAll(fs.workDir()); err != nil {
		return fmt.Errorf("failed to remove overlayfs work dir (%s): %s", fs.workDir(), err)
	}

	return nil
}
