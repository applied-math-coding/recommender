package services

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Example usage:
// model.GET("progress", func(c *gin.Context) {
// 	for progress := range services.Progress {
// 		e := services.SendSSE("message", progress, c.Writer)
// 		services.HandleError(errors.Wrap(e, "event.Render failed"), c)
// 	}
// })
// Note: some client-side dev-proxies don't correctly handle these connections.
func SendSSE(event string, data interface{}, w gin.ResponseWriter) error {
	writeContentType(w)
	w.WriteString("event:")
	w.WriteString(event)
	w.WriteString("\n")
	w.WriteString("data:")
	d, e := json.Marshal(data)
	if e != nil {
		return errors.Wrap(e, "json.Marshal failed")
	}
	w.Write(d)
	w.WriteString("\n")
	w.WriteString("\n\n")
	w.Flush()
	return nil
}

func writeContentType(w gin.ResponseWriter) {
	header := w.Header()
	header["Content-Type"] = []string{"text/event-stream"}

	if _, exist := header["Cache-Control"]; !exist {
		header["Cache-Control"] = []string{"no-cache"}
	}
}
