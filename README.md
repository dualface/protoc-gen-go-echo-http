# protoc-gen-go-echo-http

Generating Go code for Echo Web Framework from protobuf.

~

## Install

```bash
go install github.com/dualface/protoc-gen-go-echo-http/cmd/protoc-gen-go-echo-http@latest
```

~

## Example

The full example is in the `example` directory.

~

## Define the HTTP API and generate code

Create the file `auth.proto`:

```protobuf
syntax = "proto3";
package api;

option go_package = "example/api";

// AuthService is a service for authentication.
service AuthService {
  // ApplyToken returns a token for the given username and password.
  rpc ApplyToken(ApplyTokenRequest) returns (ApplyTokenResponse) {}
}

// ApplyTokenRequest is a request for ApplyToken.
message ApplyTokenRequest {
  string username = 1; // username for authentication.
  string password = 2; // password for authentication.
}

// ApplyTokenResponse is a response for ApplyToken.
message ApplyTokenResponse {
  string token = 1; // token for the given username and password.
}
```

Use the `protoc` command to generate Go code:

```bash
protoc --go_out=. --go-echo-http_out=. auth.proto
```

When the command is successfully executed,
the following files will be generated:

- `auth.pb.go`: The protobuf definition.
- `auth_http.pb.go`: The HTTP interface definition for Echo framework.

~

## Implementing the HTTP API

After generating the interface definitions,
these interfaces need to be implemented.

Create the file `auth_impl.go`:

```go
type (
	// AuthServiceImplement is a implementation of AuthService.
	AuthServiceImplement struct {
		UnimplementedAuthService
	}
)

func (s *AuthServiceImplement) ApplyToken(c echo.Context, req *ApplyTokenRequest) (resp *ApplyTokenResponse, err error) {
	username := strings.ToLower(strings.TrimSpace(req.Username))
	password := strings.TrimSpace(req.Password)

	if len(username) == 0 {
		return nil, c.String(http.StatusBadRequest, "username is required")
	}
	if len(password) == 0 {
		return nil, c.String(http.StatusBadRequest, "password is required")
	}
	if username != "testuser" || password != "testpwd" {
		return nil, c.String(http.StatusUnauthorized, "invalid username or password")
	}

	return &ApplyTokenResponse{
		Token: "you_got_a_token",
	}, nil
}
```

~

## Using the HTTP API

Using interfaces in `main.go`:

```go
func main() {
	e := echo.New()
	group := e.Group("/v1")

    // BindAuthService binds the AuthServiceImplement to the group.
	api.BindAuthService(group, &api.AuthServiceImplement{})

	e.Logger.Fatal(e.Start(":12345"))
}
```

## Testing

Start service:

```text
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.10.2
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:12345
```

Testing the interface using `curl`:

```bash
curl -v \
  'http://localhost:12345/v1/auth/apply_token' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"username":"testuser","password":"testpwd"}'
```

The expected result should be:

```text
< HTTP/1.1 200 OK
{
  "token": "you_got_a_token"
}
```

Invalid parameters can be passed in to test the error handling of the interface:

```bash
curl -v \
  'http://localhost:12345/v1/auth/apply_token' \
  -H 'Accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"username":"no","password":"no"}'
```

```text
< HTTP/1.1 401 Unauthorized
invalid username or password
```

You can keep trying.

~

## Generate Swagger documents

Install `swaggo`:

````bash

https://github.com/swaggo/swag

```bash
go install github.com/swaggo/swag/cmd/swag@latest
````

Generate swagger documents:

```bash
swag init
```

```text
> Generate swagger docs....
> Generate general API Info, search dir:./
> Generating api.ApplyTokenRequest
> Generating api.ApplyTokenResponse
> warning: route POST /auth/apply_token is declared multiple times
> create docs.go at docs/docs.go
> create swagger.json at docs/swagger.json
> create swagger.yaml at docs/swagger.yaml
```

Swagger documentation can be added to the HTTP service.
Modify the file `main.go`:

```go
import (
	"example/api"
	"example/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {
	e := echo.New()
	e.Debug = true
	e.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	e.Logger.SetLevel(log.DEBUG)

	group := e.Group("/v1")
	api.BindAuthService(group, &api.AuthServiceImplement{})

	docs.SwaggerInfo.BasePath = "/v1"
	e.GET("/v1/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":12345"))
}
```

Use your browser to view Swagger documentation:

`http://localhost:12345/v1/swagger/index.html`


![generated_swagger_docs.jpg](generated_swagger_docs.jpg)

The full example is in the `example` directory.
