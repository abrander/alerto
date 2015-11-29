package api

import (
	"sync"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/abrander/alerto/agent"
	"github.com/abrander/alerto/monitor"
)

func Run(wg sync.WaitGroup) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(static.Serve("/", static.LocalFile("/home/abrander/gocode/src/github.com/abrander/alerto/web/", true)))

	a := router.Group("/agent")
	{
		a.GET("/", func(c *gin.Context) {
			c.JSON(200, agent.AvailableAgents())
		})
	}

	m := router.Group("/monitor")
	{

		m.GET("/:id", func(c *gin.Context) {
			id := c.Param("id")

			mon, err := monitor.GetMonitor(id)
			if err == monitor.ErrorInvalidId {
				c.AbortWithError(400, err)
			} else if err != nil {
				c.AbortWithError(404, err)
			} else {
				c.JSON(200, mon)
			}
		})

		m.PUT("/:id", func(c *gin.Context) {
			var mon monitor.Monitor
			c.Bind(&mon)
			err := monitor.UpdateMonitor(&mon)
			if err != nil {
				c.AbortWithError(500, err)
			} else {
				c.JSON(200, mon)
			}
		})

		m.DELETE("/:id", func(c *gin.Context) {
			id := c.Param("id")

			err := monitor.DeleteMonitor(id)
			if err != nil {
				c.AbortWithError(500, err)
			} else {
				c.JSON(200, nil)
			}
		})

		m.POST("/new", func(c *gin.Context) {
			var mon monitor.Monitor
			c.Bind(&mon)
			err := monitor.AddMonitor(&mon)
			if err != nil {
				c.AbortWithError(500, err)
			} else {
				c.JSON(200, mon)
			}
		})

		m.GET("/", func(c *gin.Context) {
			c.JSON(200, monitor.GetAllMonitors())
		})
	}
	router.Run(":9901")

	wg.Done()
}
