// Copyright 2023 dualface. All rights reserved.
// Use of this source code is governed by a Apache2 license
// that can be found in the LICENSE file.

// Package echohttp is a generator of the Echo Web Framework.
package echohttp

import (
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FileDescriptorProto.package field number
const fileDescriptorProtoPackageFieldNumber = 2

// FileDescriptorProto.syntax field number
const fileDescriptorProtoSyntaxFieldNumber = 12

// GenerateFile generates a _http.pb.go file containing Echo scaffolding code.
func GenerateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}
	filename := file.GeneratedFilenamePrefix + "_http.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)

	// Attach all comments associated with the syntax field.
	genLeadingComments(g, file.Desc.SourceLocations().ByPath(protoreflect.SourcePath{fileDescriptorProtoSyntaxFieldNumber}))
	g.P(`// Code generated by protoc-gen-go-echo-http. DO NOT EDIT.`)
	g.P(`// versions:`)
	g.P(`// - protoc-gen-go-echo-http v`, Version)
	g.P(`// - protoc                  `, protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P(`// `, file.Desc.Path(), ` is a deprecated file.`)
	} else {
		g.P(`// source: `, file.Desc.Path())
	}
	g.P()

	// Attach all comments associated with the package field.
	genLeadingComments(g, file.Desc.SourceLocations().ByPath(protoreflect.SourcePath{fileDescriptorProtoPackageFieldNumber}))
	g.P(`package `, file.GoPackageName)
	g.P()

	r := parseFile(file)
	tmpl.Execute(g, r)

	g.P()
	return g
}

func genLeadingComments(g *protogen.GeneratedFile, loc protoreflect.SourceLocation) {
	for _, s := range loc.LeadingDetachedComments {
		g.P(protogen.Comments(s))
		g.P()
	}
	if s := loc.LeadingComments; s != "" {
		g.P(protogen.Comments(s))
		g.P()
	}
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}
