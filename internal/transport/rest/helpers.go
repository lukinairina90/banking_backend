package rest

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lukinairina90/banking_backend/internal/domain"
)

// getUserIDFromContext takes the user ID from the gin.Context.
func getUserIDFromContext(ctx *gin.Context) (int, error) {
	val, ok := ctx.Get(ctxUserIDKey)
	if !ok {
		return 0, errors.New("request context doesn't contains user id")
	}

	return val.(int), nil
}

// buildOrderingMessage takes the filter from the url and returns domain.Orderings.
func buildOrderingMessage(input string, supportedFields []string) (domain.Orderings, error) {
	if input == "" {
		return nil, nil
	}

	orderings := make(domain.Orderings)

	fieldsCondition := strings.Split(input, "|")
	for _, fieldCondition := range fieldsCondition {
		parts := strings.Split(fieldCondition, ":")
		if len(parts) != 2 {
			return nil, errors.New("wrong ordering filter format")
		}

		var isSupportedField bool
		for _, field := range supportedFields {
			if parts[0] == field {
				isSupportedField = true
			}
		}

		if !isSupportedField {
			return nil, fmt.Errorf("using unsupported ordering filter param [%s]", parts[0])
		}

		var direction string
		switch parts[1] {
		case "asc":
			direction = "asc"
		case "desc":
			direction = "desc"
		default:
			return nil, errors.New("unsupported filtering direction")
		}

		orderings[parts[0]] = direction
	}

	return orderings, nil
}

//func getUserRoleIDFromContext(ctx *gin.Context) (int, error) {
//	val, ok := ctx.Get(ctxUserRoleIDKey)
//	if !ok {
//		return 0, errors.New("request context doesn't contains user role id")
//	}
//
//	return val.(int), nil
//}
