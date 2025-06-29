package milogo_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/manuelarte/milogo"
	"github.com/manuelarte/milogo/pkg/config"
)

func TestEchoRoute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		body     any
		fields   string
		expected string
	}{
		"no query param fields": {
			body:     map[string]any{"name": "Manuel", "age": 99},
			fields:   "",
			expected: `{"name":"Manuel","age":99}`,
		},
		"query param fields, 1/2": {
			body:     map[string]any{"name": "Manuel", "age": 99},
			fields:   "name",
			expected: `{"name":"Manuel"}`,
		},
		"query param fields, 2/2": {
			body:     map[string]any{"name": "Manuel", "age": 99},
			fields:   "name,age",
			expected: `{"name":"Manuel","age":99}`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			router := setupRouter()
			url := "/echo"
			router.POST(url, func(c *gin.Context) {
				var body map[string]any

				err := c.BindJSON(&body)
				if err != nil {
					c.Status(400)

					return
				}

				c.JSON(200, body)
			})

			w := httptest.NewRecorder()

			out, err := json.Marshal(test.body)
			if err != nil {
				t.Fatal(err)
			}

			if test.fields != "" {
				url += "?fields=" + test.fields
			}

			ctx := context.Background()
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(out))
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			assert.JSONEq(t, test.expected, w.Body.String())
		})
	}
}

func TestEchoCustomHeadersRoute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		body               any
		fields             string
		customHeaders      map[string]string
		expectedHTTPStatus int
		expectedBody       string
	}{
		"query param fields, 1/2, one custom header": {
			body:               map[string]any{"name": "Manuel", "age": 99},
			fields:             "name",
			customHeaders:      map[string]string{"X-Milogo": "one_deleted"},
			expectedHTTPStatus: http.StatusOK,
			expectedBody:       `{"name":"Manuel"}`,
		},
		"query param fields, 2/2, two custom headers": {
			body:               map[string]any{"name": "Manuel", "age": 99},
			fields:             "name,age",
			customHeaders:      map[string]string{"X-Milogo": "one_deleted", "X-Trace-id": "1"},
			expectedHTTPStatus: http.StatusAccepted,
			expectedBody:       `{"name":"Manuel","age":99}`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			router := setupRouter()
			url := "/echo-headers"
			router.POST(url, func(c *gin.Context) {
				for key, value := range test.customHeaders {
					c.Writer.Header().Add(key, value)
				}

				var body map[string]any

				err := c.BindJSON(&body)
				if err != nil {
					c.Status(400)

					return
				}

				c.JSON(test.expectedHTTPStatus, body)
			})

			w := httptest.NewRecorder()

			out, err := json.Marshal(test.body)
			if err != nil {
				t.Fatal(err)
			}

			if test.fields != "" {
				url += "?fields=" + test.fields
			}

			ctx := context.Background()
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(out))
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedHTTPStatus, w.Code)

			for key, value := range test.customHeaders {
				assert.Equal(t, value, w.Header().Get(key))
			}

			assert.JSONEq(t, test.expectedBody, w.Body.String())
		})
	}
}

func TestArrayRoute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		body     map[string]any
		fields   string
		expected string
	}{
		"no query param fields": {
			body:     map[string]any{"name": "Manuel", "age": 99},
			fields:   "",
			expected: `[{"name":"Manuel","age":99}]`,
		},
		"query param one field": {
			body:     map[string]any{"name": "Manuel", "age": 99},
			fields:   "name",
			expected: `[{"name":"Manuel"}]`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			router := setupRouter()
			router.POST("/array-echo", func(c *gin.Context) {
				var body []*map[string]any

				body = append(body, &test.body)
				c.JSON(200, body)
			})

			w := httptest.NewRecorder()

			out, err := json.Marshal(test.body)
			if err != nil {
				t.Fatal(err)
			}

			url := "/array-echo"
			if test.fields != "" {
				url += "?fields=" + test.fields
			}

			ctx := context.Background()
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(out))
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			assert.JSONEq(t, test.expected, w.Body.String())
		})
	}
}

func TestEchoWrapRoute(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		body     any
		fields   string
		expected string
	}{
		"no query param fields": {
			body:     map[string]any{"name": "Manuel", "age": 99},
			fields:   "",
			expected: `{"data":{"name":"Manuel","age":99}}`,
		},
		"query param fields, 1/2": {
			body:     map[string]any{"name": "Manuel", "age": 99},
			fields:   "name",
			expected: `{"data":{"name":"Manuel"}}`,
		},
		"query param fields, 2/2": {
			body:     map[string]any{"name": "Manuel", "age": 99},
			fields:   "name,age",
			expected: `{"data":{"name":"Manuel","age":99}}`,
		},
		"query param in array, 1/2": {
			body:     []map[string]any{{"name": "Manuel", "age": 99}},
			fields:   "name,age",
			expected: `{"data":[{"name":"Manuel","age":99}]}`,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			milogoOption, _ := config.WithWrapField("data")
			router := setupRouter(milogoOption)
			url := "/echo-wrap"
			router.POST(url, func(c *gin.Context) {
				type Response struct {
					Data any `json:"data"`
				}
				c.JSON(200, Response{Data: test.body})
			})

			w := httptest.NewRecorder()

			out, err := json.Marshal(test.body)
			if err != nil {
				t.Fatal(err)
			}

			if test.fields != "" {
				url += "?fields=" + test.fields
			}

			ctx := context.Background()
			req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(out))
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code)
			assert.JSONEq(t, test.expected, w.Body.String())
		})
	}
}

func setupRouter(configOptions ...config.Option) *gin.Engine {
	r := gin.Default()
	r.Use(milogo.Milogo(configOptions...))

	return r
}
