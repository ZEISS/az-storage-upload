package main

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"

	"template/internal/cfg"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/sethvargo/go-githubactions"
	"github.com/zeiss/pkg/utilx"
)

func main() {
	ctx := context.Background()

	action := githubactions.New()

	cfg, err := cfg.NewFromInput(action)
	if err != nil {
		githubactions.Fatalf("error: %s", err)
	}

	credentials, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		githubactions.Fatalf("error: %s", err)
	}

	client, err := azblob.NewClient(cfg.AccountURL, credentials, nil)
	if err != nil {
		githubactions.Fatalf("error: %s", err)
	}

	err = filepath.WalkDir(cfg.Path, func(path string, d fs.DirEntry, errr error) error {
		root := filepath.Clean(cfg.Path)
		path = filepath.Clean(path)

		if utilx.Or(root == path, d.IsDir()) {
			return nil
		}

		file, err := os.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return err
		}
		defer file.Close()

		p, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		githubactions.Infof("uploading %s", p)

		_, err = client.UploadFile(ctx, cfg.ContainerName, p, file, nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		githubactions.Fatalf("error: %s", err)
	}

	githubactions.SetOutput("success", "true")
	githubactions.Infof("upload successful")
}
