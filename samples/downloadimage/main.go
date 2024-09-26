package main

import (
	"fmt"
	"path/filepath"

	"github.com/bougou/go-container-utils"
	ctutils "github.com/bougou/go-container-utils"
)

func main() {
	fmt.Println("download 1")
	if err := ctutils.DownloadImageTarFile("ubuntu:22.04", "linux/arm64", "/tmp"); err != nil {
		panic(err)
	}
	fmt.Println("download 2")

	if err := ctutils.DownloadImageTarFile("ubuntu:22.04", "linux/amd64", "/tmp"); err != nil {
		panic(err)
	}

	fmt.Println("download 3")

	if err := ctutils.DownloadImageTarFile("registry.cn-hangzhou.aliyuncs.com/openbayes_common/cert-manager-controller:v1.15.1", "linux/amd64", "/tmp"); err != nil {
		panic(err)
	}
	fmt.Println("load 4")

	dc := container.NewDockerContainer("docker://1234")

	tarFile, _ := ctutils.SafeImageFileName("ubuntu:22.04", "linux/amd64")
	imageTarFilePath := filepath.Join("/tmp", tarFile)
	if err := dc.LoadImage(imageTarFilePath); err != nil {
		panic(err)
	}

}
