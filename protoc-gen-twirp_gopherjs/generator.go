// Copyright 2018 Twitch Interactive, Inc.  All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the License is
// located at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// or in the "license" file accompanying this file. This file is distributed on
// an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/shutej/protobuf/proto"
	"github.com/shutej/protobuf/protoc-gen-go/descriptor"
	plugin "github.com/shutej/protobuf/protoc-gen-go/plugin"
	"github.com/shutej/twirp/internal/gen"
	"github.com/shutej/twirp/internal/gen/stringutils"
	"github.com/shutej/twirp/internal/gen/typemap"
)

type twirp struct {
	filesHandled   int
	currentPackage string // Go name of current package we're working on

	reg *typemap.Registry

	// Map to record whether we've built each package
	pkgs          map[string]string
	pkgNamesInUse map[string]bool

	// Package naming:
	genPkgName          string // Name of the package that we're generating
	fileToGoPackageName map[*descriptor.FileDescriptorProto]string

	// List of files that were inputs to the generator. We need to hold this in
	// the struct so we can write a header for the file that lists its inputs.
	genFiles []*descriptor.FileDescriptorProto

	// Output buffer that holds the bytes we want to write out for a single file.
	// Gets reset after working on a file.
	output *bytes.Buffer
}

func newGenerator() *twirp {
	t := &twirp{
		pkgs:                make(map[string]string),
		pkgNamesInUse:       make(map[string]bool),
		fileToGoPackageName: make(map[*descriptor.FileDescriptorProto]string),
		output:              bytes.NewBuffer(nil),
	}

	return t
}

func (t *twirp) Generate(in *plugin.CodeGeneratorRequest) *plugin.CodeGeneratorResponse {
	t.genFiles = gen.FilesToGenerate(in)

	// Collect information on types.
	t.reg = typemap.New(in.ProtoFile)

	// Register names of packages that we import.
	t.registerPackageName("utils")

	// Time to figure out package names of objects defined in protobuf. First,
	// we'll figure out the name for the package we're generating.
	genPkgName, err := deduceGenPkgName(t.genFiles)
	if err != nil {
		gen.Fail(err.Error())
	}
	t.genPkgName = genPkgName

	// Next, we need to pick names for all the files that are dependencies.
	for _, f := range in.ProtoFile {
		if fileDescSliceContains(t.genFiles, f) {
			// This is a file we are generating. It gets the shared package name.
			t.fileToGoPackageName[f] = t.genPkgName
		} else {
			// This is a dependency. Use its package name.
			name := f.GetPackage()
			if name == "" {
				name = stringutils.BaseName(f.GetName())
			}
			name = stringutils.CleanIdentifier(name)
			t.fileToGoPackageName[f] = name
			t.registerPackageName(name)
		}
	}

	// Showtime! Generate the response.
	resp := new(plugin.CodeGeneratorResponse)
	for _, f := range t.genFiles {
		respFile := t.generate(f)
		if respFile != nil {
			resp.File = append(resp.File, respFile)
		}
	}
	return resp
}

func (t *twirp) registerPackageName(name string) (alias string) {
	alias = name
	i := 1
	for t.pkgNamesInUse[alias] {
		alias = name + strconv.Itoa(i)
		i++
	}
	t.pkgNamesInUse[alias] = true
	t.pkgs[name] = alias
	return alias
}

// deduceGenPkgName figures out the go package name to use for generated code.
// Will try to use the explicit go_package setting in a file (if set, must be
// consistent in all files). If no files have go_package set, then use the
// protobuf package name (must be consistent in all files)
func deduceGenPkgName(genFiles []*descriptor.FileDescriptorProto) (string, error) {
	var genPkgName string
	for _, f := range genFiles {
		name, explicit := goPackageName(f)
		if explicit {
			name = stringutils.CleanIdentifier(name)
			if genPkgName != "" && genPkgName != name {
				// Make sure they're all set consistently.
				return "", errors.Errorf("files have conflicting go_package settings, must be the same: %q and %q", genPkgName, name)
			}
			genPkgName = name
		}
	}
	if genPkgName != "" {
		return genPkgName, nil
	}

	// If there is no explicit setting, then check the implicit package name
	// (derived from the protobuf package name) of the files and make sure it's
	// consistent.
	for _, f := range genFiles {
		name, _ := goPackageName(f)
		name = stringutils.CleanIdentifier(name)
		if genPkgName != "" && genPkgName != name {
			return "", errors.Errorf("files have conflicting package names, must be the same or overridden with go_package: %q and %q", genPkgName, name)
		}
		genPkgName = name
	}

	// All the files have the same name, so we're good.
	return genPkgName, nil
}

func (t *twirp) generate(file *descriptor.FileDescriptorProto) *plugin.CodeGeneratorResponse_File {
	resp := new(plugin.CodeGeneratorResponse_File)
	if len(file.Service) == 0 {
		return nil
	}

	t.generateFileHeader(file)

	t.generateImports(file)

	// For each service, generate client stubs and server
	for i, service := range file.Service {
		t.generateService(file, service, i)
	}

	resp.Name = proto.String(goFileName(file))
	resp.Content = proto.String(t.formattedOutput())
	t.output.Reset()

	t.filesHandled++
	return resp
}

func (t *twirp) generateFileHeader(file *descriptor.FileDescriptorProto) {
	t.P("//+build js")
	t.P()
	t.P("// Code generated by protoc-gen-twirp_gopherjs ", gen.Version, ", DO NOT EDIT.")
	t.P("// source: ", file.GetName())
	t.P()
	if t.filesHandled == 0 {
		t.P("/*")
		t.P("Package ", t.genPkgName, " is a generated twirp stub package.")
		t.P("This code was generated with github.com/shutej/twirp/protoc-gen-twirp_gopherjs ", gen.Version, ".")
		t.P()
		comment, err := t.reg.FileComments(file)
		if err == nil && comment.Leading != "" {
			for _, line := range strings.Split(comment.Leading, "\n") {
				line = strings.TrimPrefix(line, " ")
				// ensure we don't escape from the block comment
				line = strings.Replace(line, "*/", "* /", -1)
				t.P(line)
			}
			t.P()
		}
		t.P("It is generated from these files:")
		for _, f := range t.genFiles {
			t.P("\t", f.GetName())
		}
		t.P("*/")
	}
	t.P(`package `, t.genPkgName)
	t.P()
}

func (t *twirp) generateImports(file *descriptor.FileDescriptorProto) {
	if len(file.Service) == 0 {
		return
	}
	t.P(`import `, t.pkgs["utils"], ` "github.com/shutej/twirp/gopherjs/utils"`)
	t.P()

	// It's legal to import a message and use it as an input or output for a
	// method. Make sure to import the package of any such message. First, dedupe
	// them.
	deps := make(map[string]string) // Map of package name to quoted import path.
	ourImportPath := path.Dir(goFileName(file))
	for _, s := range file.Service {
		for _, m := range s.Method {
			defs := []*typemap.MessageDefinition{
				t.reg.MethodInputDefinition(m),
				t.reg.MethodOutputDefinition(m),
			}
			for _, def := range defs {
				importPath := path.Dir(goFileName(def.File))
				if importPath != ourImportPath {
					pkg := t.goPackageName(def.File)
					deps[pkg] = strconv.Quote(importPath)
				}
			}
		}
	}
	for pkg, importPath := range deps {
		t.P(`import `, pkg, ` `, importPath)
	}
	if len(deps) > 0 {
		t.P()
	}
}

// P forwards to g.gen.P, which prints output.
func (t *twirp) P(args ...string) {
	for _, v := range args {
		t.output.WriteString(v)
	}
	t.output.WriteByte('\n')
}

// Big header comments to makes it easier to visually parse a generated file.
func (t *twirp) sectionComment(sectionTitle string) {
	t.P()
	t.P(`// `, strings.Repeat("=", len(sectionTitle)))
	t.P(`// `, sectionTitle)
	t.P(`// `, strings.Repeat("=", len(sectionTitle)))
	t.P()
}

func (t *twirp) generateService(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto, index int) {
	servName := serviceName(service)

	t.sectionComment(servName + ` Interface`)
	t.generateTwirpInterface(file, service)

	t.sectionComment(servName + ` Protobuf Client`)
	t.generateClient("Protobuf", file, service)

	t.sectionComment(servName + ` JSON Client`)
	t.generateClient("JSON", file, service)
}

func (t *twirp) generateTwirpInterface(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	servName := serviceName(service)

	pathPrefixConst := servName + "PathPrefix"
	t.P(`// `, pathPrefixConst, ` is used for all URL paths on a twirp `, servName, ` server.`)
	t.P(`// Requests are always: POST `, pathPrefixConst, `/method`)
	t.P(`// It can be used in an HTTP mux to route twirp requests along with non-twirp requests on other routes.`)
	t.P(`const `, pathPrefixConst, ` = `, strconv.Quote(pathPrefix(file, service)))
	t.P()

	comments, err := t.reg.ServiceComments(file, service)
	if err == nil {
		t.printComments(comments)
	}
	t.P(`type `, servName, ` interface {`)
	for _, method := range service.Method {
		comments, err = t.reg.MethodComments(file, service, method)
		if err == nil {
			t.printComments(comments)
		}
		t.P(t.generateSignature(method))
		t.P()
	}
	t.P(`}`)
}

func (t *twirp) generateSignature(method *descriptor.MethodDescriptorProto) string {
	methName := methodName(method)
	inputType := t.goTypeName(method.GetInputType())
	outputType := t.goTypeName(method.GetOutputType())
	return fmt.Sprintf(`	%s(*%s) (*%s, error)`, methName, inputType, outputType)
}

// valid names: 'JSON', 'Protobuf'
func (t *twirp) generateClient(name string, file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) {
	servName := serviceName(service)
	pathPrefixConst := servName + "PathPrefix"
	structName := unexported(servName) + name + "Client"
	newClientFunc := "New" + servName + name + "Client"

	methCnt := strconv.Itoa(len(service.Method))
	t.P(`type `, structName, ` struct {`)
	t.P(`  client `, t.pkgs["utils"], `.XHRClient`)
	t.P(`  urls   [`, methCnt, `]string`)
	t.P(`}`)
	t.P()
	t.P(`// `, newClientFunc, ` creates a `, name, ` client that implements the `, servName, ` interface.`)
	t.P(`// It communicates using `, name, ` and can be configured with a custom XHRClient.`)
	t.P(`func `, newClientFunc, `(addr string, client `, t.pkgs["utils"], `.XHRClient) `, servName, ` {`)
	t.P(`  prefix := `, t.pkgs["utils"], `.URLBase(addr) + `, pathPrefixConst)
	t.P(`  urls := [`, methCnt, `]string{`)
	for _, method := range service.Method {
		t.P(`    	prefix + "`, methodName(method), `",`)
	}
	t.P(`  }`)
	t.P(`  return &`, structName, `{`)
	t.P(`    client: client,`)
	t.P(`    urls:   urls,`)
	t.P(`  }`)
	t.P(`}`)
	t.P()

	for i, method := range service.Method {
		methName := methodName(method)
		inputType := t.goTypeName(method.GetInputType())
		outputType := t.goTypeName(method.GetOutputType())

		t.P(`func (c *`, structName, `) `, methName, `(in *`, inputType, `) (*`, outputType, `, error) {`)
		t.P(`  out := new(`, outputType, `)`)
		t.P(`  err := `, t.pkgs["utils"], `.Do`, name, `Request(c.client, c.urls[`, strconv.Itoa(i), `], in, out)`)
		t.P(`  return out, err`)
		t.P(`}`)
		t.P()
	}
}

// pathPrefix returns the base path for all methods handled by a particular
// service. It includes a trailing slash. (for example
// "/twirp/twitch.example.Haberdasher/").
func pathPrefix(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) string {
	return fmt.Sprintf("/twirp/%s/", fullServiceName(file, service))
}

// pathFor returns the complete path for requests to a particular method on a
// particular service.
func pathFor(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto, method *descriptor.MethodDescriptorProto) string {
	return pathPrefix(file, service) + stringutils.CamelCase(method.GetName())
}

func (t *twirp) printComments(comments typemap.DefinitionComments) bool {
	text := strings.TrimSuffix(comments.Leading, "\n")
	if len(strings.TrimSpace(text)) == 0 {
		return false
	}
	split := strings.Split(text, "\n")
	for _, line := range split {
		t.P("// ", strings.TrimPrefix(line, " "))
	}
	return len(split) > 0
}

// Given a protobuf name for a Message, return the Go name we will use for that
// type, including its package prefix.
func (t *twirp) goTypeName(protoName string) string {
	def := t.reg.MessageDefinition(protoName)
	if def == nil {
		gen.Fail("could not find message for", protoName)
	}

	var prefix string
	if pkg := t.goPackageName(def.File); pkg != t.genPkgName {
		prefix = pkg + "."
	}

	var name string
	for _, parent := range def.Lineage() {
		name += parent.Descriptor.GetName() + "_"
	}
	name += def.Descriptor.GetName()
	return prefix + name
}

func (t *twirp) goPackageName(file *descriptor.FileDescriptorProto) string {
	return t.fileToGoPackageName[file]
}

func (t *twirp) formattedOutput() string {
	// Reformat generated code.
	fset := token.NewFileSet()
	raw := t.output.Bytes()
	ast, err := parser.ParseFile(fset, "", raw, parser.ParseComments)
	if err != nil {
		// Print out the bad code with line numbers.
		// This should never happen in practice, but it can while changing generated code,
		// so consider this a debugging aid.
		var src bytes.Buffer
		s := bufio.NewScanner(bytes.NewReader(raw))
		for line := 1; s.Scan(); line++ {
			fmt.Fprintf(&src, "%5d\t%s\n", line, s.Bytes())
		}
		gen.Fail("bad Go source code was generated:", err.Error(), "\n"+src.String())
	}

	out := bytes.NewBuffer(nil)
	err = (&printer.Config{Mode: printer.TabIndent | printer.UseSpaces, Tabwidth: 8}).Fprint(out, fset, ast)
	if err != nil {
		gen.Fail("generated Go source code could not be reformatted:", err.Error())
	}

	return out.String()
}

func unexported(s string) string { return strings.ToLower(s[:1]) + s[1:] }

func fullServiceName(file *descriptor.FileDescriptorProto, service *descriptor.ServiceDescriptorProto) string {
	name := stringutils.CamelCase(service.GetName())
	if pkg := pkgName(file); pkg != "" {
		name = pkg + "." + name
	}
	return name
}

func pkgName(file *descriptor.FileDescriptorProto) string {
	return file.GetPackage()
}

func serviceName(service *descriptor.ServiceDescriptorProto) string {
	return stringutils.CamelCase(service.GetName())
}

func serviceStruct(service *descriptor.ServiceDescriptorProto) string {
	return unexported(serviceName(service)) + "Server"
}

func methodName(method *descriptor.MethodDescriptorProto) string {
	return stringutils.CamelCase(method.GetName())
}

func fileDescSliceContains(slice []*descriptor.FileDescriptorProto, f *descriptor.FileDescriptorProto) bool {
	for _, sf := range slice {
		if f == sf {
			return true
		}
	}
	return false
}
