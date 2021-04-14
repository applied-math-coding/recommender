package main

import (
	"main/routes"
	"main/services"
	"runtime"

	"github.com/gin-gonic/gin"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	services.InitDb(false)

	server := gin.Default()
	server.Static("/app", "./client/build")
	server.NoRoute(func(c *gin.Context) {
		c.File("./client/build/index.html")
	})
	api := server.Group("/api")
	routes.DataRoute(api)
	routes.ModelRoute(api)

	server.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
