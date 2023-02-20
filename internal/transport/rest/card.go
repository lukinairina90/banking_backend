package rest

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/banking_backend/internal/transport/rest/messages"
)

// Card transport layer struct.
type Card struct {
	cardService CardService
}

// NewCard constructor for Card.
func NewCard(cardService CardService) *Card {
	return &Card{cardService: cardService}
}

// InjectRoutes injects routes to global router.
func (t Card) InjectRoutes(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	accountCard := r.Group("/account/:id/card").Use(middlewares...)
	{
		accountCard.POST("/", t.createCard)
		accountCard.GET("/", t.getCardListByAccount)
		accountCard.GET("/:card_id", t.getCard)
	}

	card := r.Group("/card").Use(middlewares...)
	{
		card.GET("/", t.GetCardListByUser)
	}
}

// createCard gin handler function for creation card endpoint.
// [POST] /account/:id/card
func (t Card) createCard(ctx *gin.Context) {
	accountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	domainCard, err := t.cardService.CreateCard(ctx, accountID, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("card creation error", err))
		return
	}

	messageCard := messages.Card{
		Id:             domainCard.Id,
		AccountID:      domainCard.AccountID,
		CardNumber:     domainCard.CardNumber,
		CardholderName: domainCard.CardholderName,
		ExpirationDate: domainCard.ExpirationDate,
		CvvCode:        domainCard.CvvCode,
	}

	ctx.JSON(http.StatusCreated, messageCard)
}

// GetCardListByUser gin handler function for get card list for user endpoint.
// [GET] /card
func (t Card) GetCardListByUser(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	domainListCards, err := t.cardService.GetCardListUser(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting cards list error", err))
		return
	}

	messageListCards := make([]messages.Card, 0)
	for _, card := range domainListCards {
		messageCard := messages.Card{
			Id:             card.Id,
			AccountID:      card.AccountID,
			CardNumber:     card.CardNumber,
			CardholderName: card.CardholderName,
			ExpirationDate: card.ExpirationDate,
			CvvCode:        card.CvvCode,
		}
		messageListCards = append(messageListCards, messageCard)
	}

	ctx.JSON(http.StatusOK, messageListCards)
}

// getCardListByAccount gin handler function for get card list for account endpoint.
// [GET] /account/:id/card
func (t Card) getCardListByAccount(ctx *gin.Context) {
	accountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	domainListCards, err := t.cardService.GetCardListByAccount(ctx, userID, accountID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting cards list for account error", err))
		return
	}

	messageListCards := make([]messages.Card, 0)
	for _, card := range domainListCards {
		messageCard := messages.Card{
			Id:             card.Id,
			AccountID:      card.AccountID,
			CardNumber:     card.CardNumber,
			CardholderName: card.CardholderName,
			ExpirationDate: card.ExpirationDate,
			CvvCode:        card.CvvCode,
		}
		messageListCards = append(messageListCards, messageCard)
	}

	ctx.JSON(http.StatusOK, messageListCards)
}

// getCard gin handler function for get card endpoint.
// [GET] /account/:id/card/:card_id
func (t Card) getCard(ctx *gin.Context) {
	accountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	cardID, err := strconv.Atoi(ctx.Param("card_id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"card_id\" request param", err))
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	domainCard, err := t.cardService.GetCard(ctx, cardID, accountID, userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewNotFoundError("card not found", err))
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting card error", err))
		}

		return
	}

	messageCard := messages.Card{
		Id:             domainCard.Id,
		AccountID:      domainCard.AccountID,
		CardNumber:     domainCard.CardNumber,
		CardholderName: domainCard.CardholderName,
		ExpirationDate: domainCard.ExpirationDate,
		CvvCode:        domainCard.CvvCode,
	}

	ctx.JSON(http.StatusOK, messageCard)
}
