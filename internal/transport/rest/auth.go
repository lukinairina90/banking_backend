package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/banking_backend/internal/transport/rest/messages"
)

// Auth transport layer struct.
type Auth struct {
	userService UserService
}

// NewAuth constructor for Auth.
func NewAuth(userService UserService) *Auth {
	return &Auth{userService: userService}
}

// InjectRoutes injects routes to global router.
func (a *Auth) InjectRoutes(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	auth := r.Group("/auth").Use(middlewares...)
	{
		auth.POST("/sign-up", a.signUp)
		auth.POST("/sign-in", a.signIn)
		auth.GET("/refresh", a.refresh)
	}

	user := r.Group("/user").Use(middlewares...)
	{
		user.POST("/:id/block", a.blockUser)
		user.POST("/:id/unblock", a.unblockUser)
	}
}

// signUp gin handler function for user registration endpoint.
// [POST] /auth/sign-up
func (a *Auth) signUp(ctx *gin.Context) {
	var inp messages.SignUpInput

	if err := ctx.ShouldBindJSON(&inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("validation request body error", err))
		return
	}

	domainInp := inp.ToDomain()

	err := a.userService.SignUp(ctx, domainInp)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("signing up error", err))
		return
	}

	ctx.JSON(http.StatusOK, inp)
}

// signIn gin handler function for user login endpoint.
// [POST] /auth/sign-in
func (a *Auth) signIn(ctx *gin.Context) {
	var inp messages.SignInInput
	if err := ctx.ShouldBindJSON(&inp); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("validation request body error", err))
		return
	}

	domainInp := inp.ToDomain()

	accessToken, refreshToken, err := a.userService.SignIn(ctx, domainInp)
	if err != nil {
		if errors.Is(err, messages.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, NewNotFoundError("user not found", err))
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("signing in error", err))
		return
	}

	ctx.SetCookie("refresh-token", refreshToken, 3600, "/auth", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
}

// refresh gin handler function for refresh the token endpoint.
// [GET] /auth/refresh
func (a *Auth) refresh(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh-token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("get cookie from request error", err))
		return
	}

	accessToken, refreshToken, err := a.userService.RefreshTokens(ctx, cookie)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("refresh token error", err))
		return
	}

	ctx.SetCookie("refresh-token", refreshToken, 3600, "/auth", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"token": accessToken,
	})
}

// blockUser gin handler function to block a user endpoint.
// [GET] /user/:id/block
func (a *Auth) blockUser(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting current user error", err))
		return
	}

	blockUserID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	err = a.userService.BlockUser(ctx, blockUserID, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("blocking user error", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// unblockUser gin handler function to unblock a user endpoint.
// [GET] /user/:id/unblock
func (a *Auth) unblockUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewBadRequestError("wrong \"id\" request param", err))
		return
	}

	err = a.userService.UnblockUser(ctx, userID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("unblocking user error", err))
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
