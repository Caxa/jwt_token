package handler

import (
	"jwt_token/pkg"
	"jwt_token/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// signUp обрабатывает регистрацию пользователя.
func (h *Handler) signUp(c *gin.Context) {
	var input pkg.User
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// signInInput структура запроса для входа.
type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// signIn обрабатывает вход пользователя.
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	authService, ok := h.services.Authorization.(*service.AuthService)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "service casting error")
		return
	}

	id, err := authService.AuthenticateUser(input.Email, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "signed in successfully",
		"user_id": id,
	})
}
