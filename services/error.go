package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func HandleError(e error, c *gin.Context, doPanic bool) {
	if e != nil {
		if c != nil {
			c.String(http.StatusInternalServerError, "Some error happened!")
		}
		if errors.Cause(e) != nil {
			fmt.Printf("%+v\n", e)
		} else {
			fmt.Println(e)
		}
		if doPanic {
			panic(e)
		}
	}
}
