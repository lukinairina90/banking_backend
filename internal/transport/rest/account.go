package rest

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/banking_backend/internal/domain"
	"github.com/lukinairina90/banking_backend/internal/transport/rest/messages"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Account transport layer struct.
type Account struct {
	accService AccountService
}

// NewAccount constructor for Account.
func NewAccount(accService AccountService) *Account {
	return &Account{accService: accService}
}

// InjectRoutes injects routes to global router.
func (t Account) InjectRoutes(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	accounts := r.Group("/account").Use(middlewares...)
	{
		accounts.POST("/", t.createAccount)
		accounts.GET("/", t.getAccountsList)
		accounts.GET("/:id", t.getAccount)
		accounts.DELETE("/:id", t.deleteAccount)
		accounts.POST("/:id/deposit", t.depositAccount)
		accounts.POST("/:id/transfer", t.transferAccount)
		accounts.POST("/:id/block", t.blockAccount)
		accounts.POST("/:id/unblock", t.unblockAccount)
	}
}

// createAccount gin handler function for account creation endpoint.
// [POST] /account/
func (t Account) createAccount(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	var req messages.CreateAccountRequestBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("request body validation error", err))
		return
	}

	domainAccount, err := t.accService.Create(ctx, userID, req.CurrencyID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("account creation error", err))
		return
	}

	messageAccount := messages.Account{
		ID:         domainAccount.ID,
		Iban:       domainAccount.Iban,
		UserID:     domainAccount.UserID,
		CurrencyID: domainAccount.CurrencyID,
		Blocked:    domainAccount.Blocked,
		Amount:     domainAccount.Amount,
	}

	ctx.JSON(http.StatusCreated, messageAccount)
}

// getAccountsList gin handler function for get account list endpoint.
// [GET] /account/
func (t Account) getAccountsList(ctx *gin.Context) {
	// http://localhost:8080/account?order=amount:desc|id:asc&page=1&per-page=100
	// http://localhost:8080/account?page=1&per-page=100

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
			ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"page\" query param", err))
			return
		}

		pag.Page = page
	}

	if sPerPage != "" {
		perPage, err := strconv.Atoi(sPerPage)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"per-page\" query param", err))
			return
		}

		pag.PerPage = perPage
	}

	orderings, err := buildOrderingMessage(rawOrdering, []string{"id", "iban", "amount"})
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong ordering query param", err))
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	domainAccountsList, err := t.accService.GetAccountsList(ctx, userID, pag, orderings)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting account list error", err))
		return
	}

	messageAccountsList := make([]messages.Account, 0)
	for _, account := range domainAccountsList {
		messageAccount := messages.Account{
			ID:         account.ID,
			Iban:       account.Iban,
			UserID:     account.UserID,
			CurrencyID: account.CurrencyID,
			Blocked:    account.Blocked,
			Amount:     account.Amount,
		}
		messageAccountsList = append(messageAccountsList, messageAccount)
	}
	ctx.JSON(http.StatusOK, messageAccountsList)
}

// getAccount gin handler function for get account endpoint.
// [GET] /account/:id
func (t Account) getAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		logrus.WithError(err).Error("getting user id from context error")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	domainAccount, err := t.accService.GetAccount(ctx, id, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewNotFoundError("account not found", err))
			return
		default:
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting account error", err))
			return
		}
	}

	messageAccount := messages.Account{
		ID:         domainAccount.ID,
		Iban:       domainAccount.Iban,
		UserID:     domainAccount.UserID,
		CurrencyID: domainAccount.CurrencyID,
		Blocked:    domainAccount.Blocked,
		Amount:     domainAccount.Amount,
	}

	ctx.JSON(http.StatusOK, messageAccount)
}

// deleteAccount gin handler function for delete account endpoint.
// [DELETE] /account/:id
func (t Account) deleteAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	if err := t.accService.DeleteAccount(ctx, id); err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("deletion account error", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// depositAccount gin handler function for deposit account endpoint.
// [POST] /account/:id/deposit
func (t Account) depositAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	var req messages.DepositAccountRequestBody
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("validation request body error", err))
		return
	}

	err = t.accService.DepositAccount(ctx, id, req.Amount)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("deposit founds into account error", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// transferAccount gin handler function for transfer money account endpoint.
// [POST] /account/:id/transfer
func (t Account) transferAccount(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	var req messages.TransferAccountRequestBody
	if err = ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("validation request body error", err))
		return
	}

	err = t.accService.TransferAccount(ctx, id, userID, req.Amount, req.Iban)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("transferring founds error", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// blockAccount gin handler function for block account endpoint.
// [POST] /account/:id/block
func (t Account) blockAccount(ctx *gin.Context) {
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

	err = t.accService.BlockAccount(ctx, accountID, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("blocking user error", err))
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// blockAccount gin handler function for unblock account endpoint.
// [POST] /account/:id/unblock
func (t Account) unblockAccount(ctx *gin.Context) {
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

	err = t.accService.UnblockAccount(ctx, accountID, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("unblocking account error", err))
	}

	ctx.JSON(http.StatusNoContent, nil)
}
