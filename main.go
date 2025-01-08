package main

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"template/internal/cfg"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/h2non/filetype"
	"github.com/sethvargo/go-githubactions"
	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/utilx"
)

// GetContentType ...
func GetContentType(seeker io.ReadSeeker) (string, error) {
	// At most the first 512 bytes of data are used:
	// https://golang.org/src/net/http/sniff.go?s=646:688#L11
	buff := make([]byte, 512)

	_, err := seeker.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	bytesRead, err := seeker.Read(buff)
	if utilx.NotEmpty(err) && !errors.Is(err, io.EOF) {
		return "", err
	}

	// Slice to remove fill-up zero values which cause a wrong content type detection in the next step
	buff = buff[:bytesRead]

	kind, err := filetype.Match(buff)
	if err != nil && !errors.Is(err, filetype.ErrEmptyBuffer) {
		return "", err
	}

	if errors.Is(err, filetype.ErrEmptyBuffer) {
		return "application/octet-stream", nil
	}

	return kind.MIME.Type, nil
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

		kind, err := GetContentType(file)
		if err != nil {
			return err
		}

		opts := &azblob.UploadFileOptions{
			HTTPHeaders: &blob.HTTPHeaders{
				BlobContentType: cast.Ptr(kind),
			},
		}

		githubactions.Infof("uploading %s (%s)", p, kind)

		_, err = client.UploadFile(ctx, cfg.ContainerName, p, file, opts)
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
