package rest

import (
	"net/http"
	"strings"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ctxUserIDKey     = "user-id"
	ctxUserRoleIDKey = "user-role-id"
)

const AuthorizationHeaderName = "Authorization"

// LoggingMiddleware middleware for logging
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fields := logrus.Fields{
			"method":          c.Request.Method,
			"uri":             c.Request.RequestURI,
			"request-in-time": t.Format(time.RFC3339),
		}

		c.Next()

		dur := time.Since(t)
		fields["request-handling-duration"] = dur.Milliseconds()

		logrus.WithFields(fields).Info()
	}
}

// AuthMiddleware middleware for api, takes a token from the request, checks the authorization token, checks if the user is blocked, sets the user id and role id to the context.
func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewUnauthorizedError("token is missing or invalid", err))
			return
		}

		userID, roleID, err := a.userService.ParseToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, NewUnauthorizedError("wrong authorization token", err))
			return
		}

		checkBlockUser, err := a.userService.CheckBlockUser(c, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, NewNotFoundError("user block check error user not found", err))
			return
		}

		if checkBlockUser {
			c.AbortWithStatusJSON(http.StatusForbidden, NewForbiddenError("user is blocked", err))
			return
		}

		c.Set(ctxUserIDKey, userID)
		c.Set(ctxUserRoleIDKey, roleID)

		c.Next()
	}
}

// RBACMiddleware checks if the user has access to certain endpoints.
func RBACMiddleware(enforcer casbin.IEnforcer, roleRepository RoleRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var group, method, path string
		r, exists := c.Get(ctxUserRoleIDKey)
		if !exists {
			group = "anonymous"
		} else {
			roleID := r.(int)
			role, err := roleRepository.GetByID(c.Request.Context(), roleID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("getting role error", err))
				return
			}

			group = role.Name
		}

		method = c.Request.Method
		path = c.Request.URL.Path

		ok, err := enforcer.Enforce(group, path, method)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, NewInternalServerError("checking permissions error", err))
			return
		}

		if ok {
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusForbidden, NewForbiddenError("user does not have rights to perform an operation", nil))
			return
		}
	}
}

// getTokenFromRequest takes a token from a request.
func getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.GetHeader(AuthorizationHeaderName)
	if header == "" {
		return "", errors.Wrap(nil, "empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.Wrap(nil, "invalid authorization header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.Wrap(nil, "token is empty")
	}

	return headerParts[1], nil
}
