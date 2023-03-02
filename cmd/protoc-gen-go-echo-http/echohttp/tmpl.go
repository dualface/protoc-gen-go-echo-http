package echohttp

import (
	"html/template"
)

var methodSignature = `%s(c echo.Context, req *%s) (resp *%s, err error)`

var tmpl, _ = template.New("tmpl").Parse(`
import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	{{- range $import := .Imports }}
	{{ $import.Alias -}} "{{ $import.Import }}"
	{{- end }}
)

{{ range $service := .Services }}

	{{ if $service.LeadingComments -}}
	// {{ $service.LeadingComments -}}
	// @Description {{ $service.LeadingComments -}}
	{{ end -}}
	type {{ $service.Name }} interface {
		{{- range $method := $service.Methods }}
		{{ if $method.LeadingComments }}// {{ $method.LeadingComments }}{{ end -}}
		{{ $method.Signature }}{{- if $method.TrailingComments }} // {{ $method.TrailingComments }}{{ end }}
		{{ end }}
		mustEmbedUnimplemented{{ $service.Name }}()
	}

	type Unimplemented{{ $service.Name }} struct {}

	{{- range $method := $service.Methods }}

	{{ if $method.LeadingComments -}}
 	// {{ $method.LeadingComments -}}
	// @Description  {{ $method.LeadingComments -}}
	{{ end -}}
	// @Tags         {{ $service.Route }}
	// @Accept       json
	// @Produce      json
	// @Param        req body {{ $method.ReqType }} true "{{ $method.ReqTypeParamComment }}"
	{{ range $an := $service.Annotations -}}
	{{- if eq $an.Name "jwt" -}}
	// @Param 		 Authorization header string true "Bearer _JWT_TOKEN_"
	{{ end -}}
	{{ end -}}
	// @Success      200 {object} {{ $method.RespType }}
	// @Router       /{{ $service.Route }}/{{ $method.Route }} [post]
	func (s *Unimplemented{{ $service.Name }}) {{ $method.Signature }} {
		return nil, c.JSON(http.StatusNotImplemented, "/{{ $service.Route }}/{{ $method.Route }} not implemented")
	}

	{{ end }}

	func (s *Unimplemented{{ $service.Name }}) mustEmbedUnimplemented{{ $service.Name }}() {
	}

	// Bind{{ $service.Name }} register http handlers
	func Bind{{ $service.Name }}(g *echo.Group, s {{ $service.Name }}) *echo.Group {
		genErrorResp := func(c echo.Context, err error) error {
			errMap := map[string]interface{}{
				"code": 500,
				"message": err.Error(),
			}
			return c.JSON(http.StatusOK, map[string]interface{}{
				"err": errMap,
			})
		}

		handlePanic := func(c echo.Context) error {
			if r := recover(); r != nil {
				log.Printf("%v\n%s", r, debug.Stack())
				return genErrorResp(c, fmt.Errorf("%v", r))
			}
			return nil
		}

		group := g.Group("/{{ $service.Route }}")

		{{- range $an := $service.Annotations }}
			{{ if eq $an.Name "jwt" }}
			// Auth with middleware JWT
			group.Use(setupJWTAuth())
			{{ end -}}
		{{ end -}}

		{{ range $index, $method := $service.Methods }}
		group.POST("/{{ $method.Route }}", func(c echo.Context) error {
			defer handlePanic(c)

			var req {{ $method.ReqType }}
			if err := c.Bind(&req); err != nil {
				return c.JSON(http.StatusBadRequest, err.Error())
			}

			resp, err := s.{{ $method.Name }}(c, &req)
			if err != nil {
				return genErrorResp(c, err)
			}
			if resp == nil {
				return nil
			}
			return c.JSON(http.StatusOK, resp)
		})
		{{ end }}

		return group
	}

{{ end }}

`)
