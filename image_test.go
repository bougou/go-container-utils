package container

import (
	"fmt"
	"testing"
)

func Test_DownloadImageTarFile(t *testing.T) {
	fmt.Println("xxx")
	if err := DownloadImageTarFile("ubuntu:22.04", "linux/amd64", "/tmp"); err != nil {
		t.Error(err)
	}
}
