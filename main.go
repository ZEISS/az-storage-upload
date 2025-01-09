package main

import (
	"context"
	"io/fs"
	"mime"
	"os"
	"path/filepath"

	"template/internal/cfg"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/sethvargo/go-githubactions"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/utilx"
)

// GetContentType ...
func GetContentType(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".htm", ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	default:
		return mime.TypeByExtension(ext)
	}
}

// nolint:gocyclo
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

		ct := GetContentType(path)

		opts := &azblob.UploadFileOptions{
			HTTPHeaders: &blob.HTTPHeaders{
				BlobContentType: cast.Ptr(ct),
			},
		}

		githubactions.Infof("uploading %s (%s)", p, ct)

		_, err = client.UploadFile(ctx, cfg.ContainerName, p, file, opts)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		githubactions.Fatalf("error: %s", err)
	}

	githubactions.Infof("upload successful")
}
