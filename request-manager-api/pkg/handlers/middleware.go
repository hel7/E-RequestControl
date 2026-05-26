package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "UserID"
)

func (h *Handlers) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		c.Abort()
		return
	}

	headerParts := strings.Fields(header)
	if len(headerParts) < 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		c.Abort()
		return
	}

	userID, roleID, err := h.service.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		c.Abort()
		return
	}

	c.Set(userCtx, userID)
	c.Set("RoleID", int(roleID))
	c.Next()
}

func getUserID(c *gin.Context) (int, error) {
	id, exists := c.Get(userCtx)
	if !exists {
		return 0, errors.New("User id not found")
	}

	userID, ok := id.(int)
	if !ok {
		return 0, errors.New("User id not valid type")
	}

	return userID, nil
}

func (h *Handlers) adminRequired(c *gin.Context) {
	roleIDValue, exists := c.Get("RoleID")
	if !exists {
		newErrorResponse(c, http.StatusUnauthorized, "User role not found")
		c.Abort()
		return
	}

	roleID, ok := roleIDValue.(int)
	if !ok {
		newErrorResponse(c, http.StatusUnauthorized, "Invalid role type")
		c.Abort()
		return
	}

	if roleID != 1 {
		newErrorResponse(c, http.StatusForbidden, "Admin role required")
		c.Abort()
		return
	}

	c.Next()
}
