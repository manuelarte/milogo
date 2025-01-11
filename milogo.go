package milogo

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/manuelarte/milogo/pkg"
	"strconv"
	"strings"
)

// customResponseWriter captures the response body for modification
type customResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the data in the body buffer
func (w *customResponseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func Milogo(configOptions ...pkg.ConfigOption) gin.HandlerFunc {
	config := pkg.DefaultConfig(configOptions...)

	return func(c *gin.Context) {
		// Create a custom response writer
		writer := &customResponseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		// Continue to the next middleware/handler
		c.Next()

		// If the content-type is JSON, modify the JSON body
		if isPartialResponseRequest(c, config) {
			var jsonData interface{}
			if err := json.Unmarshal(writer.body.Bytes(), &jsonData); err == nil {
				fields := c.Query(config.QueryParamField)

				wrappedJsonData := jsonData
				if config.WrapperField != "" {
					wrappedJsonData = jsonData.(map[string]interface{})[config.WrapperField]
				}
				if partialResponseFields, errParsing := config.Parser.Parse(fields); errParsing == nil &&
					pkg.Filter(wrappedJsonData, partialResponseFields) == nil {
					modifiedBody, err := json.Marshal(jsonData)
					if err == nil {
						c.Writer = writer.ResponseWriter // Set back to original writer
						_, _ = c.Writer.Write(modifiedBody)
						c.Header("Content-Length", strconv.Itoa(len(modifiedBody)))

						return
					}
				}
			}
		}

		// If JSON parsing fails or content-type is not applicable, write the original body
		c.Writer = writer.ResponseWriter
		_, _ = c.Writer.Write(writer.body.Bytes())
	}
}

func isPartialResponseRequest(c *gin.Context, config pkg.Config) bool {
	isJson := strings.Contains(c.Writer.Header().Get("Content-Type"), "application/json")
	isFieldQuery := c.Query(config.QueryParamField) != ""
	isNotBadStatus := c.Writer.Status() < 300 && c.Writer.Status() > 199

	return isJson && isFieldQuery && isNotBadStatus
}
