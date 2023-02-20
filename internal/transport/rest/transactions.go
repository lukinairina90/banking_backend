package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/transport/rest/messages"
)

// Transaction transport layer struct.
type Transaction struct {
	TransactionService TransactionService
}

// NewTransaction constructor for transaction.
func NewTransaction(transactionService TransactionService) *Transaction {
	return &Transaction{TransactionService: transactionService}
}

// InjectRoutes injects routes to global router.
func (t Transaction) InjectRoutes(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	transaction := r.Group("account/:id/transaction").Use(middlewares...)
	{
		transaction.GET("/", t.getTransactionList)
	}
}

// getTransactionList gin handler function for get list transaction endpoint.
// [GET] /account/:id/transaction
func (t Transaction) getTransactionList(ctx *gin.Context) {
	// http://localhost:8080/account/1/transaction?order=amount:desc|id:asc&page=1&per-page=100
	// http://localhost:8080/account/1/transaction?page=1&per-page=100

	sPage := ctx.Request.URL.Query().Get("page")
	sPerPage := ctx.Request.URL.Query().Get("per-page")
	rawOrdering := ctx.Request.URL.Query().Get("order")

	pag := domain.Paginator{
		Page:    1,
		PerPage: 5,
	}

	if sPage != "" {
		page, err := strconv.Atoi(sPage)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"page\" request param", err))
			return
		}

		pag.Page = page
	}

	if sPerPage != "" {
		perPage, err := strconv.Atoi(sPerPage)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"per-page\" request param", err))
			return
		}

		pag.PerPage = perPage
	}

	orderings, err := buildOrderingMessage(rawOrdering, []string{"id", "date_updated"})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong ordering query param", err))
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	accountID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	domainTransactionsList, err := t.TransactionService.GetTransactionList(ctx, accountID, userID, orderings, pag)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting transactions list error", err))
		return
	}

	messageTransactionsList := make([]messages.Transaction, 0, len(domainTransactionsList))
	for _, transaction := range domainTransactionsList {
		messageTransaction := messages.Transaction{
			ID:          transaction.ID,
			FromAccount: transaction.FromAccount,
			ToAccount:   transaction.ToAccount,
			Amount:      transaction.Amount,
			Type:        transaction.Type,
			Status:      transaction.Status,
			DateCreated: transaction.DateCreated,
			DateUpdated: transaction.DateUpdated,
		}
		messageTransactionsList = append(messageTransactionsList, messageTransaction)
	}

	ctx.JSON(http.StatusOK, messageTransactionsList)
}
