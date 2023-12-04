package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/pluswing/datasync/data"
	"github.com/pluswing/datasync/file"
	"github.com/spf13/cobra"
)

func Upload(target string, fileName string, conf data.StorageGcsType) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	f, err := os.Open(target)
	cobra.CheckErr(err)
	defer f.Close()

	// ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	// defer cancel()
	var uploadPath = ""
	if conf.Dir == "" {
		uploadPath = fileName
	} else {
		uploadPath = filepath.Join(conf.Dir, fileName)
	}

	o := client.Bucket(conf.Bucket).Object(uploadPath)

	// o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	_, err = io.Copy(wc, f)
	cobra.CheckErr(err)

	err = wc.Close()
	cobra.CheckErr(err)
}

func Download(target string, conf data.StorageGcsType) string {
	// TODO client は一度作ったものを使い回す。
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	var filePath = ""
	if conf.Dir == "" {
		filePath = target
	} else {
		filePath = filepath.Join(conf.Dir, target)
	}

	o := client.Bucket(conf.Bucket).Object(filePath)

	tmpDir, err := file.MakeTempFile()
	cobra.CheckErr(err)

	tmpFile := filepath.Join(tmpDir, target)
	f, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY, 0644)
	cobra.CheckErr(err)
	defer f.Close()

	rc, err := o.NewReader(ctx)
	cobra.CheckErr(err)
	defer rc.Close()

	_, err = io.Copy(f, rc)
	cobra.CheckErr(err)

	// TODO tmpFileを消す処理を忘れずに。
	return tmpFile
}

func Exists(target string, conf data.StorageGcsType) bool {
	// TODO client は一度作ったものを使い回す。
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	cobra.CheckErr(err)
	defer client.Close()

	var filePath = ""
	if conf.Dir == "" {
		filePath = target
	} else {
		filePath = filepath.Join(conf.Dir, target)
	}

	o := client.Bucket(conf.Bucket).Object(filePath)

	attrs, _ := o.Attrs(ctx)

	return attrs != nil
}
