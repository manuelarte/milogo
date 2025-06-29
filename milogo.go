package milogo

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/manuelarte/milogo/internal/parser"
	"github.com/manuelarte/milogo/pkg/config"
)

// customResponseWriter captures the response body for modification.
type customResponseWriter struct {
	gin.ResponseWriter

	body *bytes.Buffer
}

// Write captures the data in the body buffer.
func (w *customResponseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func Milogo(configOptions ...config.Option) gin.HandlerFunc {
	cfg := config.DefaultConfig(configOptions...)

	return func(c *gin.Context) {
		// Create a custom response writer
		writer := &customResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		// Continue to the next middleware/handler
		c.Next()

		// If the content-type is JSON, modify the JSON body
		//nolint:nestif // Refactor later
		if jsonBody, isPartialResponse := isPartialResponseRequest(c, cfg); isPartialResponse {
			fields := c.Query(cfg.QueryParamField)

			wrappedJSONData := jsonBody
			if cfg.WrapperField != "" {
				if wrappedField, isWrappedAJSON := jsonBody.(map[string]any); isWrappedAJSON {
					wrappedJSONData = wrappedField[cfg.WrapperField]
				} else {
					return
				}
			}

			if partialResponseFields, errParsing := cfg.Parser.Parse(fields); errParsing == nil &&
				parser.Filter(wrappedJSONData, partialResponseFields) == nil {
				modifiedBody, errMarsh := json.Marshal(jsonBody)
				if errMarsh == nil {
					c.Writer = writer.ResponseWriter // Set back to original writer
					_, _ = c.Writer.Write(modifiedBody)
					c.Header("Content-Length", strconv.Itoa(len(modifiedBody)))

					return
				}
			}
		}

		// If JSON parsing fails or content-type is not applicable, write the original body
		c.Writer = writer.ResponseWriter
		_, _ = c.Writer.Write(writer.body.Bytes())
	}
}

func isPartialResponseRequest(c *gin.Context, cfg config.Config) (any, bool) {
	is300 := 300
	is199 := 199

	isJSON := strings.Contains(c.Writer.Header().Get("Content-Type"), "application/json")
	isFieldQuery := c.Query(cfg.QueryParamField) != ""

	isNotBadStatus := c.Writer.Status() < is300 && c.Writer.Status() > is199
	if isJSON && isFieldQuery && isNotBadStatus {
		if customWriter, isCustomWriter := c.Writer.(*customResponseWriter); isCustomWriter {
			var jsonBody any

			copiedBody := make([]byte, len(customWriter.body.Bytes()))
			copy(copiedBody, customWriter.body.Bytes())

			err := json.Unmarshal(copiedBody, &jsonBody)
			if err == nil {
				return jsonBody, true
			}
		}
	}

	return nil, false
}
