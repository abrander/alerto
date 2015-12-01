package api

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/abrander/alerto/monitor"
	"github.com/abrander/alerto/plugins"
)

type (
	Message struct {
		Type    string      `json:"type"`
		Payload interface{} `json:"payload"`
	}

	Status struct {
		Uptime  time.Duration `json:"uptime"`
		Clock   time.Time     `json:"clock"`
		Started time.Time     `json:"start"`
	}
)

var (
	StartTime  time.Time
	wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func init() {
	StartTime = time.Now()
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	ticker := time.Tick(time.Second)
	changes := monitor.SubscribeChanges()

	status := Status{
		Started: StartTime,
	}

	for {
		select {
		case t := <-ticker:
			status.Clock = t
			status.Uptime = t.Sub(StartTime)
			err := conn.WriteJSON(Message{Type: "status", Payload: status})
			if err != nil {
				goto unsubscribe
			}
		case msg := <-changes:
			err := conn.WriteJSON(msg)
			if err != nil {
				goto unsubscribe
			}
		}
	}

unsubscribe:
	monitor.UnsubscribeChanges(changes)
}

func Run(wg sync.WaitGroup) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())

	router.Use(static.Serve("/", static.LocalFile("/home/abrander/gocode/src/github.com/abrander/alerto/web/", true)))

	router.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	a := router.Group("/agent")
	{
		a.GET("/", func(c *gin.Context) {
			c.JSON(200, plugins.AvailableAgents())
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
