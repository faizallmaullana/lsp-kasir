package handler

import (
	"net/http"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/http/dto"
	"faizalmaulana/lsp/http/middleware"
	"faizalmaulana/lsp/http/services"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type AuthenticationHandler struct {
	svc  services.AuthenticationService
	sess services.SessionService
	cfg  *conf.Config
}

func NewAuthenticationHandler(s services.AuthenticationService, sess services.SessionService, cfg *conf.Config) *AuthenticationHandler {
	return &AuthenticationHandler{svc: s, sess: sess, cfg: cfg}
}

func (h *AuthenticationHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/login", middleware.LoginRateLimiter(), h.login)
	rg.POST("/refresh", middleware.JWTMiddleware(h.cfg), h.refresh)
}

func (h *AuthenticationHandler) login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}

	user, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}

	session, err := h.sess.Create(user.IdUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to create session"))
		return
	}

	token, err := services.GenerateToken(h.cfg, user.IdUser, session.IdSession, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to generate token"))
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", gin.H{"user": user, "token": token}))
}

func (h *AuthenticationHandler) refresh(c *gin.Context) {
	claimsVal, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}
	claims, ok := claimsVal.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}

	userID, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	sessionID, _ := claims["session_id"].(string)
	if userID == "" || sessionID == "" {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}

	newToken, err := services.GenerateToken(h.cfg, userID, sessionID, "", role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to refresh token"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", gin.H{"token": newToken}))
}
