package routes

import (
	"main/domain"
	"main/services"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func DataRoute(api *gin.RouterGroup) {
	data := api.Group("/data")
	broadcasters := make(map[int]*services.Broadcaster)
	broadcastersMut := sync.Mutex{}

	data.POST("", func(c *gin.Context) {
		fileHeader, e1 := c.FormFile("file")
		if e1 != nil {
			services.HandleError(errors.Wrap(e1, "c.FormFile fails"), c, false)
			return
		}
		var broadcasterId int
		broadcastersMut.Lock()
		broadcasterId = len(broadcasters)
		broadcasters[broadcasterId] = services.CreateBroadcaster(
			make(chan domain.ProgressMessage))
		broadcastersMut.Unlock()
		go services.ProcessFileData(fileHeader, broadcasters[broadcasterId])
		c.JSON(http.StatusOK, gin.H{"id": broadcasterId})
	})

	data.GET("progress/:id", func(c *gin.Context) {
		broadcasterId, paramIdError := strconv.Atoi(c.Param("id"))
		if paramIdError != nil {
			services.HandleError(errors.Wrap(paramIdError, "strconv.Atoi"), c, false)
			return
		}
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			}}
		ws, e1 := upgrader.Upgrade(c.Writer, c.Request, nil)
		if e1 != nil {
			services.HandleError(errors.Wrap(e1, "upgrader.Upgrade failed"), c, false)
			return
		}
		defer ws.Close()
		broadcaster := broadcasters[broadcasterId]
		if broadcaster == nil {
			progress := domain.ProgressMessage{
				Message:         "No data import is running.",
				State:           domain.ProgressState.Error,
				ShowProgressBar: false}
			e2 := ws.WriteJSON(progress)
			if e2 != nil {
				services.HandleError(errors.Wrap(e2, "ws.WriteJSONfailed"), c, false)
				return
			}
			return
		}
		quitWs := make(chan bool)
		broadcaster.BroadcasterMut.Lock()
		if broadcaster.LastProgressMessage != nil && broadcaster.LastProgressMessage.State != domain.ProgressState.Running {
			e3 := ws.WriteJSON(broadcaster.LastProgressMessage)
			if e3 != nil {
				services.HandleError(errors.Wrap(e3, "ws.WriteJSONfailed"), c, false)
			}
			close(quitWs)
		} else {
			broadcaster.Listener = func(progress domain.ProgressMessage, done bool) {
				if done {
					quitWs <- true
				} else {
					e4 := ws.WriteJSON(progress)
					if e4 != nil {
						services.HandleError(errors.Wrap(e4, "ws.WriteJSONfailed"), c, false)
						quitWs <- true
					}
				}
			}
		}
		broadcaster.BroadcasterMut.Unlock()
		<-quitWs
	})

	data.DELETE("cancel/:id", func(c *gin.Context) {
		broadcasterId, paramIdError := strconv.Atoi(c.Param("id"))
		if paramIdError != nil {
			services.HandleError(errors.Wrap(paramIdError, "strconv.Atoi"), c, false)
			return
		}
		broadcaster := broadcasters[broadcasterId]
		if broadcaster == nil {
			c.String(http.StatusOK, "")
			return
		}
		broadcaster.Cancel = true
		c.String(http.StatusOK, "")
	})
}
