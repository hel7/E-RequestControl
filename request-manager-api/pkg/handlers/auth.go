package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	Request_Manager "request_manager_api"
)

type registrationInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handlers) register(c *gin.Context) {
	var input Request_Manager.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"UserID": id,
	})
}
func (h *Handlers) registerAdmin(c *gin.Context) {
	var input Request_Manager.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.service.Authorization.CreateAdmin(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"AdminID": id,
	})
}
func (h *Handlers) login(c *gin.Context) {
	var input registrationInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.service.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
func (h *Handlers) logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		newErrorResponse(c, http.StatusBadRequest, "Authorization header is missing")
		return
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		newErrorResponse(c, http.StatusBadRequest, "Invalid Authorization header format")
		return
	}

	token := authHeader[len(bearerPrefix):]

	if err := h.service.Authorization.InvalidateToken(token); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "Failed to invalidate token: "+err.Error())
		return
	}

	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
		"actions": []string{
			"remove_local_storage_tokens",
			"clear_application_state",
			"redirect_to_login",
		},
	})
}
