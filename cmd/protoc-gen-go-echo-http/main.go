// Copyright 2023 dualface. All rights reserved.
// Use of this source code is governed by a Apache2 license
// that can be found in the LICENSE file.

// protoc-gen-go-echo-http is a plugin to generate Go code for
// the Echo Web Framework.
//
// For more information about the Echo Web Framework, see:
// https://echo.labstack.com/
//
// Install it by building this program and making it accessible within
// your PATH with the name:
//
// protoc-gen-go-echo-http
//
// The 'go' suffix becomes part of the argument for the protocol compiler,
// such that it can be invoked as:
//
// protoc --go-echo-http_out=. path/to/file.proto
//
// This generates Echo scaffolding code for the API defined by file.proto.
// With that input, the output will be written to:
//
// path/to/file_http.pb.go
//
// For more information about the protocol-buffers, see:
// https://developers.google.com/protocol-buffers/
package main

import (
	"flag"
	"fmt"

	"github.com/dualface/protoc-gen-go-echo-http/echohttp"

	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("protoc-gen-go-echo-http %v\n", echohttp.Version)
		return
	}

	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = gengo.SupportedFeatures

		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			echohttp.GenerateFile(gen, f)
		}

		return nil
	})
}
