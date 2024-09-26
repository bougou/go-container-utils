package container

//build !windows

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

func SafeImageFileName(image string, imagePlatform string) (imageTarFileName, imageDigestFileName string) {
	safeFileName := fmt.Sprintf("%s_%s", imagePlatform, image)
	safeFileName = strings.ReplaceAll(safeFileName, "/", "_")
	safeFileName = strings.ReplaceAll(safeFileName, ":", "_")

	imageTarFileName = fmt.Sprintf("%s.tar", safeFileName)
	imageDigestFileName = fmt.Sprintf("%s.digest", safeFileName)

	return
}

func DownloadImageTarFile(image string, imagePlatform string, saveTarFileDir string) error {
	if saveTarFileDir == "" {
		return fmt.Errorf("saveTarFileDir can not be empty")
	}

	if imagePlatform == "" {
		return fmt.Errorf("imagePlatform can not be empty")
	}

	imageTarFileName, imageDigestFileName := SafeImageFileName(image, imagePlatform)
	imageTarFilePath := filepath.Join(saveTarFileDir, imageTarFileName)
	imageDigestFilePath := filepath.Join(saveTarFileDir, imageDigestFileName)

	opts := []regclient.ImageOpts{}
	p, err := platform.Parse(imagePlatform)
	if err != nil {
		return fmt.Errorf("parse platform failed, err: %s", err)
	}

	rc := regclient.New()
	ctx := context.Background()

	r, err := ref.New(image)
	if err != nil {
		return fmt.Errorf("create ref failed, err: %s", err)
	}

	opts = append(opts, regclient.ImageWithExportRef(r))

	m, err := rc.ManifestGet(ctx, r, regclient.WithManifestPlatform(p))
	if err != nil {
		return fmt.Errorf("get manifest failed, err: %s", err)
	}

	if mi, ok := m.(manifest.Imager); ok {
		d, err := mi.GetConfig()
		if err != nil {
			return fmt.Errorf("get image config failed, err: %s", err)
		}

		fmt.Println("repo digest", d.Digest.String())
	}

	imageDigest := m.GetDescriptor().Digest.String()
	fmt.Println("image digest", imageDigest)

	r = r.SetDigest(imageDigest)

	var w io.Writer
	w, err = os.Create(imageTarFilePath)
	if err != nil {
		return fmt.Errorf("create image tar file failed, err: %s", err)
	}

	if err := rc.ImageExport(ctx, r, w, opts...); err != nil {
		return fmt.Errorf("import export failed, err: %s", err)
	}

	if err := os.WriteFile(imageDigestFilePath, []byte(imageDigest+"\n"), 0644); err != nil {
		return fmt.Errorf("write image digest file failed, err: %s", err)
	}

	return nil
}
