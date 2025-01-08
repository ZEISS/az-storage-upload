package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"template/internal/cfg"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/sethvargo/go-githubactions"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/conv"
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

		fileinfo, err := file.Stat()
		if err != nil {
			return err
		}

		p, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		githubactions.Infof("uploading %s", p)

		filesize := fileinfo.Size()
		buffer := make([]byte, filesize)

		_, err = file.Read(buffer)
		if err != nil {
			return err
		}

		hash := md5.Sum(buffer)

		ct := http.DetectContentType(buffer)
		opts := &azblob.UploadFileOptions{
			Metadata: map[string]*string{
				"Content-Type":   cast.Ptr(ct),
				"Content-MD5":    cast.Ptr(hex.EncodeToString(hash[:])),
				"Content-Length": cast.Ptr(conv.String(filesize)),
			},
		}

		_, err = client.UploadBuffer(ctx, cfg.ContainerName, p, buffer, opts)
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
