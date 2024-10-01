package main

import (
	"fmt"

	ctutils "github.com/bougou/go-container-utils"
)

func download(image string, imagePlatform string, saveDir string) {
	fmt.Printf(">>> download image (%s), platform (%s) to dir (%s)\n", image, imagePlatform, saveDir)
	if err := ctutils.DownloadImageTarFile(image, imagePlatform, saveDir); err != nil {
		panic(err)
	}
}
func load(runtime ctutils.Runtime, image string, imagePlatform string, saveDir string) {
	fmt.Printf(">>> runtime (%s) load image (%s), platform (%s) from dir (%s)\n", runtime, image, imagePlatform, saveDir)

	fakeRuntimeContaienrID := fmt.Sprintf("%s://%s", runtime, "fake-container-id")

	c, err := ctutils.NewContainer(fakeRuntimeContaienrID)
	if err != nil {
		panic(err)
	}

	imageTarFilePath := ctutils.SafeImageTarFilePath(image, imagePlatform, saveDir)
	if err := c.LoadImage(imageTarFilePath); err != nil {
		panic(err)
	}
}

func main() {
	saveDir := "/tmp/image-cache"
	download("ubuntu:22.04", "linux/arm64", saveDir)
	download("ubuntu:22.04", "linux/amd64", saveDir)
	download("registry.cn-hangzhou.aliyuncs.com/openbayes_common/cert-manager-controller:v1.15.1", "linux/amd64", saveDir)
	load("docker", "ubuntu:22.04", "linux/amd64", saveDir)
}
