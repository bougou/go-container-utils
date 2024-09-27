package container

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/regclient/regclient"
	"github.com/regclient/regclient/types/manifest"
	"github.com/regclient/regclient/types/platform"
	"github.com/regclient/regclient/types/ref"
)

func SafeImageFileName(image string, imagePlatform string) (safeFileName string) {
	safeFileName = fmt.Sprintf("%s_%s", imagePlatform, image)
	safeFileName = strings.ReplaceAll(safeFileName, "/", "_")
	safeFileName = strings.ReplaceAll(safeFileName, ":", "_")
	return
}

func LoadImageTarFile(image string, imagePlatform string, saveDir string) error {
	dc := NewDockerContainer("docker://fakeid")

	safeFileName := SafeImageFileName(image, imagePlatform)
	imageFileDir := filepath.Join(saveDir, safeFileName)
	imageTarFilePath := filepath.Join(imageFileDir, safeFileName+".tar")

	if err := dc.LoadImage(imageTarFilePath); err != nil {
		return fmt.Errorf("load image failed, err: %s", err)
	}

	return nil
}

// image is the url of the image.
// imagePlatform example: linux/amd64, linux/arm64
func DownloadImageTarFile(image string, imagePlatform string, saveDir string) error {
	if saveDir == "" {
		return fmt.Errorf("saveDir can not be empty")
	}

	if imagePlatform == "" {
		return fmt.Errorf("imagePlatform can not be empty")
	}

	safeFileName := SafeImageFileName(image, imagePlatform)
	imageTarFileName := fmt.Sprintf("%s.tar", safeFileName)
	imageIDFileName := fmt.Sprintf("%s.id", safeFileName)

	imageFileDir := filepath.Join(saveDir, safeFileName)
	imageTarFilePath := filepath.Join(imageFileDir, imageTarFileName)
	imageIDFilePath := filepath.Join(imageFileDir, imageIDFileName)

	if err := os.MkdirAll(imageFileDir, 0700); err != nil {
		return fmt.Errorf("create image file dir failed, err: %s", err)
	}

	rc := regclient.New()
	ctx := context.Background()

	r, err := ref.New(image)
	if err != nil {
		return fmt.Errorf("create ref failed, err: %s", err)
	}

	p, err := platform.Parse(imagePlatform)
	if err != nil {
		return fmt.Errorf("parse platform failed, err: %s", err)
	}

	m, err := rc.ManifestGet(ctx, r, regclient.WithManifestPlatform(p))
	if err != nil {
		return fmt.Errorf("get manifest failed, err: %s", err)
	}

	imageDigest := m.GetDescriptor().Digest.String()
	if imageDigest == "" {
		return fmt.Errorf("got empty image digest")
	}
	// this can make sure we reference a unique image
	// because it will unset the tag, so we need to set the tag later.
	r = r.SetDigest(imageDigest)

	mi, ok := m.(manifest.Imager)
	if !ok {
		return fmt.Errorf("manifest does not implement Imager interface")
	}
	d, err := mi.GetConfig()
	if err != nil {
		return fmt.Errorf("get image config failed, err: %s", err)
	}
	imageID := d.Digest.String()
	if imageID == "" {
		return fmt.Errorf("got empty image id")
	}

	var w io.Writer
	w, err = os.Create(imageTarFilePath)
	if err != nil {
		return fmt.Errorf("create image tar file failed, err: %s", err)
	}

	opts := []regclient.ImageOpts{}
	eRef, err := ref.New(image)
	if err != nil {
		return fmt.Errorf("cannot parse %s: %w", image, err)
	}
	// reset the tag here.
	opts = append(opts, regclient.ImageWithExportRef(eRef))

	if err := rc.ImageExport(ctx, r, w, opts...); err != nil {
		return fmt.Errorf("import export failed, err: %s", err)
	}

	if err := os.WriteFile(imageIDFilePath, []byte(imageID+"\n"), 0644); err != nil {
		return fmt.Errorf("write image digest file failed, err: %s", err)
	}

	return nil
}
