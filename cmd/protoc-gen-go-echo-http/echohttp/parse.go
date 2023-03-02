package echohttp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dualface/protoc-gen-go-echo-http/echohttp/annotation"
	"github.com/dualface/protoc-gen-go-echo-http/echohttp/extimport"

	"google.golang.org/protobuf/compiler/protogen"
)

type (
	// File Save the analysis results of the file
	File struct {
		Imports  extimport.Imports
		Services []Service
	}

	// Service Used to populate the service interface definition in the template
	Service struct {
		Name             string
		Route            string
		LeadingComments  string
		TrailingComments string
		Methods          []Method
		Annotations      []annotation.Annotation
	}

	// Method Used to populate the service method in the template
	Method struct {
		ServiceName         string
		Name                string
		Route               string
		Signature           string
		LeadingComments     string
		TrailingComments    string
		ReqType             string
		ReqTypeParamComment string
		RespType            string
		Annotations         []annotation.Annotation
		SecurityType        string
	}
)

var (
	// removeServiceSuffixReg Remove the last 'Service' word from the service name
	removeServiceSuffixReg = regexp.MustCompile("(.+)Service$")
	// addUnderscoreReg Add an underscore to each uppercase letter
	addUnderscoreReg = regexp.MustCompile("([A-Z]+)")
)

// Analyze the structure of the protobuf file to generate the data structure used to populate the template
func parseFile(file *protogen.File) *File {
	importer := extimport.New(file)
	result := &File{
		Services: make([]Service, 0, len(file.Services)),
		Imports:  importer.Imports,
	}

	for _, service := range file.Services {
		sv := Service{
			Name:             service.GoName,
			LeadingComments:  strings.TrimLeft(service.Comments.Leading.String(), "/ "),
			TrailingComments: strings.TrimLeft(service.Comments.Trailing.String(), "/ "),
			Route:            strings.ToLower(addUnderscore(removeServiceSuffixReg.ReplaceAllString(service.GoName, "$1"))),
			Annotations:      annotation.ParseAnnotations(service.Comments),
		}
		checkAnnotations(sv.Annotations, importer)
		security := checkSecurity(sv.Annotations)

		for _, method := range service.Methods {
			pm := fetchFirstLineComment(method.Input.Comments.Leading)
			if len(pm) == 0 {
				pm = method.Input.GoIdent.GoName
			}

			inputGoName, _ := importer.AddFromMessage(method.Input)
			outputGoName, _ := importer.AddFromMessage(method.Output)

			m := Method{
				ServiceName:         sv.Name,
				Name:                method.GoName,
				LeadingComments:     strings.TrimLeft(method.Comments.Leading.String(), "/ "),
				TrailingComments:    strings.TrimLeft(method.Comments.Trailing.String(), "/ "),
				Route:               strings.ToLower(addUnderscore(method.GoName)),
				ReqType:             inputGoName,
				ReqTypeParamComment: pm,
				RespType:            outputGoName,
				Annotations:         annotation.ParseAnnotations(method.Comments),
			}
			checkAnnotations(m.Annotations, importer)

			if len(security) > 0 {
				m.SecurityType = security
			}

			m.Signature = fmt.Sprintf(methodSignature, m.Name, m.ReqType, m.RespType)
			sv.Methods = append(sv.Methods, m)
		}

		result.Services = append(result.Services, sv)
	}

	return result
}

func checkAnnotations(as []annotation.Annotation, i *extimport.Importer) {
	for _, a := range as {
		switch a.Name {
		}
	}
}

func checkSecurity(as []annotation.Annotation) string {
	for _, a := range as {
		switch a.Name {
		case "jwt":
			return "jwt"
		}
	}
	return ""
}

func fetchFirstLineComment(comments protogen.Comments) string {
	parts := strings.Split(string(comments), "\n")
	if len(parts) > 0 {
		return strings.Trim(parts[0], " \t\n/")
	}
	return ""
}

func addUnderscore(s string) string {
	return strings.TrimLeft(addUnderscoreReg.ReplaceAllString(s, "_$1"), "_")
}
