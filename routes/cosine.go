package routes

import (
	"main/domain"
	"main/services"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func CosineRoute(model *gin.RouterGroup) {
	cosine := model.Group("/cosine")
	broadcasters := make(map[int]*services.Broadcaster)
	broadcastersMut := sync.Mutex{}

	cosine.POST("create", func(c *gin.Context) {
		modelParams := ModelParams{}
		e1 := c.ShouldBindJSON(&modelParams)
		if e1 != nil {
			services.HandleError(errors.Wrap(e1, "c.ShouldBindJSON failed"), c, false)
			return
		}
		var broadcasterId int
		broadcastersMut.Lock()
		broadcasterId = len(broadcasters)
		broadcasters[broadcasterId] = services.CreateBroadcaster(
			make(chan domain.ProgressMessage))
		broadcastersMut.Unlock()
		go services.ExtractAndSaveCosines(services.ItemSets, modelParams.Support, broadcasters[broadcasterId])
		c.JSON(http.StatusOK, gin.H{"id": broadcasterId})
	})

	cosine.GET("progress/:id", services.CreateProgressEndpointHandle(broadcasters))

	cosine.DELETE("cancel/:id", services.CreateCancelHandle(broadcasters))

	cosine.POST("apply", func(c *gin.Context) {
		items := make(domain.ItemSet, 0)
		e1 := c.ShouldBindJSON(&items)
		if e1 != nil {
			services.HandleError(errors.Wrap(e1, "c.ShouldBindJSON failed"), c, false)
			return
		}
		recommendations, e2 := services.MakeCosineRecommendation(items)
		if e2 != nil {
			services.HandleError(errors.WithStack(e2), c, false)
			return
		}
		c.JSON(http.StatusOK, recommendations)
	})

	cosine.GET("stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, services.ComputeCosineStatistics())
	})
}
