package api

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type (
	// AuthServiceImplement is a implementation of AuthService.
	AuthServiceImplement struct {
		UnimplementedAuthService
	}
)

// ApplyToken returns a token for the given username and password.
// @Description  ApplyToken returns a token for the given username and password.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        req body ApplyTokenRequest true "ApplyTokenRequest is a request for ApplyToken."
// @Success      200 {object} ApplyTokenResponse
// @Router       /auth/apply_token [post]
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
