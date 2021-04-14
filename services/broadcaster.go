package services

import (
	"main/domain"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

func CreateBroadcaster(progressChan chan domain.ProgressMessage) *Broadcaster {
	b := &Broadcaster{
		ProgressChan:   progressChan,
		Cancel:         false,
		BroadcasterMut: sync.Mutex{}}
	b.Run()
	return b
}

/*
Broadcaster broadcasts emitions on progress to connections.
It closes the connection when progress get closed.
*/
type Broadcaster struct {
	ProgressChan        chan domain.ProgressMessage
	Cancel              bool
	LastProgressMessage *domain.ProgressMessage
	BroadcasterMut      sync.Mutex
	Listener            func(progress domain.ProgressMessage, done bool)
}

func (b *Broadcaster) Run() {
	go func() {
		for progress := range b.ProgressChan {
			b.BroadcasterMut.Lock()
			b.LastProgressMessage = &progress
			if b.Listener != nil {
				b.Listener(progress, false)
			}
			b.BroadcasterMut.Unlock()
		}
		b.BroadcasterMut.Lock()
		if b.Listener != nil {
			b.Listener(*b.LastProgressMessage, true)
		}
		b.BroadcasterMut.Unlock()
	}()
}

func CreateCancelHandle(broadcasters map[int]*Broadcaster) gin.HandlerFunc {
	return func(c *gin.Context) {
		broadcasterId, paramIdError := strconv.Atoi(c.Param("id"))
		if paramIdError != nil {
			HandleError(errors.Wrap(paramIdError, "strconv.Atoi"), c, false)
			return
		}
		broadcaster := broadcasters[broadcasterId]
		if broadcaster == nil {
			c.String(http.StatusOK, "")
			return
		}
		broadcaster.Cancel = true
		c.String(http.StatusOK, "")
	}
}

func CreateProgressEndpointHandle(broadcasters map[int]*Broadcaster) gin.HandlerFunc {
	return func(c *gin.Context) {
		broadcasterId, paramIdError := strconv.Atoi(c.Param("id"))
		if paramIdError != nil {
			HandleError(errors.Wrap(paramIdError, "strconv.Atoi"), c, false)
			return
		}
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			}}
		ws, e1 := upgrader.Upgrade(c.Writer, c.Request, nil)
		if e1 != nil {
			HandleError(errors.Wrap(e1, "upgrader.Upgrade failed"), c, false)
			return
		}
		defer ws.Close()
		broadcaster := broadcasters[broadcasterId]
		if broadcaster == nil {
			progress := domain.ProgressMessage{
				Message:         "No model is running",
				State:           domain.ProgressState.Error,
				ShowProgressBar: false}
			e2 := ws.WriteJSON(progress)
			if e2 != nil {
				HandleError(errors.Wrap(e2, "ws.WriteJSONfailed"), c, false)
				return
			}
			return
		}
		quitWs := make(chan bool)
		broadcaster.BroadcasterMut.Lock()
		if broadcaster.LastProgressMessage != nil && broadcaster.LastProgressMessage.State != domain.ProgressState.Running {
			e3 := ws.WriteJSON(broadcaster.LastProgressMessage)
			if e3 != nil {
				HandleError(errors.Wrap(e3, "ws.WriteJSONfailed"), c, false)
			}
			close(quitWs)
		} else {
			broadcaster.Listener = func(progress domain.ProgressMessage, done bool) {
				if done {
					quitWs <- true
				} else {
					e4 := ws.WriteJSON(progress)
					if e4 != nil {
						HandleError(errors.Wrap(e4, "ws.WriteJSONfailed"), c, false)
						quitWs <- true
					}
				}
			}
		}
		broadcaster.BroadcasterMut.Unlock()
		<-quitWs
	}
}
