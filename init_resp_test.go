package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	e := echo.New()
	e.GET("/ok", func(c echo.Context) error {
		return jsonSuccess(c, echo.Map{
			"message": "",
		})
	})

	a, _ := request("GET", "/ok", e)
	assert.Equal(t, a, http.StatusOK)
}

func request(method, path string, e *echo.Echo) (int, string) {
	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}
