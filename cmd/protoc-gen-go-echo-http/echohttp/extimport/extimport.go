// Package extimport Handling of libraries to be imported from external sources
package extimport

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type (
	// Importer Handling of libraries to be imported from external sources
	Importer struct {
		Imports Imports

		fileImport string
	}

	// Import Libraries that need to be imported from outside
	Import struct {
		Import string // 导入路径
		Alias  string // 包别名
		GoName string // 类型名称
	}

	// Imports External Import List
	Imports map[string]Import
)

// New Create a new Importer
func New(f *protogen.File) *Importer {
	return &Importer{
		fileImport: strings.Trim(f.GoImportPath.String(), "/\""),
		Imports:    make(map[string]Import),
	}
}

// AddFromMessage Add a message type to the import list
func (x *Importer) AddFromMessage(m *protogen.Message) (string, bool) {
	goName := m.GoIdent.GoName
	path := strings.Trim(m.GoIdent.GoImportPath.String(), "/\"")
	if strings.Compare(path, x.fileImport) == 0 {
		return goName, false
	}

	if _, ok := x.Imports[goName]; ok {
		return goName, false
	}

	alias := path
	if p := strings.LastIndex(alias, "/"); p > 0 {
		alias = alias[p+1:]
	}
	x.Imports[goName] = Import{
		Import: path,
		Alias:  alias,
		GoName: goName,
	}

	return alias + "." + goName, true
}

// Add Add a type to the import list
func (x *Importer) Add(goName, path string) {
	if _, ok := x.Imports[goName]; ok {
		return
	}

	x.Imports[goName] = Import{
		Import: path,
		GoName: goName,
	}
}
