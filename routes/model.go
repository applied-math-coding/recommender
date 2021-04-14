package routes

import (
	"github.com/gin-gonic/gin"
)

type ModelParams = struct {
	Support int `json:"support"`
}

func ModelRoute(api *gin.RouterGroup) {
	model := api.Group("/model")
	AprioriRoute(model)
	CosineRoute(model)
}
