/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/chainguard-dev/clog"
	_ "github.com/chainguard-dev/clog/gcp/init"
	"github.com/chainguard-dev/terraform-infra-common/pkg/rotate"
	"github.com/kelseyhightower/envconfig"
)

type rotateConfig struct {
	Bucket        string        `envconfig:"BUCKET" required:"true"`
	FlushInterval time.Duration `envconfig:"FLUSH_INTERVAL" default:"3m"`
	LogPath       string        `envconfig:"LOG_PATH" required:"true"`
}

func main() {
	var rc rotateConfig
	if err := envconfig.Process("", &rc); err != nil {
		clog.Fatalf("Error processing environment: %v", err)
	}

	// log the filenames and file sizes of the files in the log path
	clog.Infof("Log path: %s", rc.LogPath)
	files, err := os.ReadDir(rc.LogPath)
	if err != nil {
		clog.Fatalf("Failed to read directory: %v", err)

	}
	// todo try this
	clog.Infof("Files in log path: %d", len(files))
	//for _, file := range files {
	//	fi, err := file.Info()
	//	if err != nil {
	//		clog.Errorf("Failed to get file info: %v", err)
	//		continue
	//	}
	//	clog.Infof("File: %s, Size: %d", file.Name(), fi.Size())
	//}

	uploader := rotate.NewUploader(rc.LogPath, rc.Bucket, rc.FlushInterval)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := uploader.Run(ctx); err != nil {
		clog.Fatalf("Failed to run the uploader: %v", err)
	}
}
