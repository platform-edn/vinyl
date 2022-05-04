//go:build mage
// +build mage

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/sh"
)

var (
	baseDir = getMageDir()
)

func getMageDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return dir
}

// updates grpc boilerplate
func Proto() error {
	protopath := filepath.Join(baseDir, "proto")
	// get files in proto path
	files, err := ioutil.ReadDir(protopath)
	if err != nil {
		return fmt.Errorf("could not get files in %s: %s", baseDir, err)
	}

	// get the generated protobuf files for each proto file for go
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".proto") {

			err = sh.Run("protoc", "--proto_path="+protopath, "--go-grpc_out=.", file.Name())
			if err != nil {
				return fmt.Errorf("could not create go proto files: %s", err)
			}

			err = sh.Run("protoc", "--proto_path="+protopath, "--go_out=.", file.Name())
			if err != nil {
				return fmt.Errorf("could not create go proto files: %s", err)
			}
		}
	}

	return nil
}
