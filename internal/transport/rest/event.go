package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/banking_backend/internal/transport/rest/messages"
)

// Event transport layer struct.
type Event struct {
	eventService EventService
}

// NewEvent constructor for Event.
func NewEvent(eventService EventService) *Event {
	return &Event{eventService: eventService}
}

// InjectRoutes injects routes to global router.
func (t Event) InjectRoutes(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	events := r.Group("/event").Use(middlewares...)
	{
		events.GET("/", t.getEventList)
	}
}

// getEventList gin handler function for get list event endpoint.
// [GET] /event/
func (t Event) getEventList(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	domainList, err := t.eventService.GetEventList(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting events list error", err))
		return
	}

	var list []messages.Event
	for _, event := range domainList {
		list = append(list, messages.Event{
			ID:       event.ID,
			UserID:   event.UserID,
			Type:     event.Type.String(),
			Message:  event.Message,
			Metadata: event.Metadata,
			DateTime: event.DateTime,
		})
	}

	ctx.JSON(http.StatusOK, list)
}
