package routes

import (
	"main/domain"
	"main/services"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func AprioriRoute(model *gin.RouterGroup) {
	apriori := model.Group("/apriori")
	broadcasters := make(map[int]*services.Broadcaster)
	broadcastersMut := sync.Mutex{}

	apriori.POST("create", func(c *gin.Context) {
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
		go services.ExtractAndSaveRules(
			services.ItemSets,
			modelParams.Support,
			broadcasters[broadcasterId])
		c.JSON(http.StatusOK, gin.H{"id": broadcasterId})
	})

	apriori.GET("progress/:id", services.CreateProgressEndpointHandle(broadcasters))

	apriori.DELETE("cancel/:id", services.CreateCancelHandle(broadcasters))

	apriori.GET("stats", func(c *gin.Context) {
		ruleStats, e := services.ComputeRuleStatistics()
		if e != nil {
			services.HandleError(errors.WithStack(e), c, false)
		}
		c.JSON(http.StatusOK, ruleStats)
	})

	apriori.POST("apply", func(c *gin.Context) {
		items := make(domain.ItemSet, 0)
		e1 := c.ShouldBindJSON(&items)
		if e1 != nil {
			services.HandleError(errors.Wrap(e1, "c.ShouldBindJSON failed"), c, false)
			return
		}
		recommendations, e2 := services.MakeAprioriRecommendation(items)
		if e2 != nil {
			services.HandleError(errors.WithStack(e2), c, false)
			return
		}
		c.JSON(http.StatusOK, recommendations)
	})

	apriori.GET("examples", func(c *gin.Context) {
		rules, e := services.FindExampleRules()
		if e != nil {
			services.HandleError(errors.WithStack(e), c, false)
			return
		}
		c.JSON(http.StatusOK, rules)
	})
}
