package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IndalAwalaikal/warung-pos/backend/utils"
	"github.com/gin-gonic/gin"
)

type NotificationController struct{}

func NewNotificationController() *NotificationController {
	return &NotificationController{}
}

// Stream opens an SSE stream to push notifications to the client
func (c *NotificationController) Stream(ctx *gin.Context) {
	w := ctx.Writer
	r := ctx.Request

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	ch := make(chan string)
	utils.NotifierInstance.AddClient(ch)
	defer utils.NotifierInstance.RemoveClient(ch)

	// send a welcome event
	fmt.Fprintf(w, "data: %s\n\n", "connected")
	flusher.Flush()

	notify := r.Context().Done()
	for {
		select {
		case <-notify:
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			// send as SSE data
			fmt.Fprintf(w, "data: %s\n\n", msg)
			flusher.Flush()
		// periodic ping to keep connection alive
		case <-time.After(30 * time.Second):
			fmt.Fprintf(w, "data: %s\n\n", "ping")
			flusher.Flush()
		}
	}
}
