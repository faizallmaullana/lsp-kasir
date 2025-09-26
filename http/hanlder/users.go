package handler

import (
	"fmt"
	"net/http"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/http/dto"
	"faizalmaulana/lsp/http/middleware"
	"faizalmaulana/lsp/http/services"
	"faizalmaulana/lsp/models/entity"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

type UsersHandler struct {
	cfg     *conf.Config
	profile services.ProfilesService
	Users   services.UsersService
}

func NewUsersHandler(cfg *conf.Config, profile services.ProfilesService, users services.UsersService) *UsersHandler {
	return &UsersHandler{cfg: cfg, profile: profile, Users: users}
}

func (h *UsersHandler) Register(rr *gin.RouterGroup) {
	rg := rr.Group("/profile")

	rg.GET("/me", middleware.JWTMiddleware(h.cfg), h.me)
	rg.POST("", middleware.JWTMiddleware(h.cfg), h.createProfile)
	rg.PUT("/:id", middleware.JWTMiddleware(h.cfg), h.updateProfile)
	rg.DELETE("/:id", middleware.JWTMiddleware(h.cfg), h.deleteProfile)
	rg.PUT("/email", middleware.JWTMiddleware(h.cfg), h.updateEmail)
}

func (h *UsersHandler) getUserIDFromClaims(c *gin.Context) (string, bool) {
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return "", false
	}
	claims, ok := claimsVal.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return "", false
	}
	subRaw, ok := claims["sub"].(string)
	if !ok || subRaw == "" {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return "", false
	}
	return subRaw, true
}

func (h *UsersHandler) me(c *gin.Context) {
	userID, ok := h.getUserIDFromClaims(c)
	if !ok {
		return
	}

	user, err := h.Users.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var profileSummaries []dto.ProfileSummary
	if user.Profiles != nil {
		for _, p := range user.Profiles {
			profileSummaries = append(profileSummaries, dto.ProfileSummary{
				IdProfile: p.IdProfile,
				Name:      p.Name,
				Contact:   p.Contact,
				Address:   p.Address,
				ImageUrl:  p.ImageUrl,
			})
		}
	}

	resp := dto.MeResponse{UserID: user.IdUser, Email: user.Email, Role: user.Role, Profiles: profileSummaries}
	fmt.Println("Me claims extracted for user:", user.IdUser)
	c.JSON(http.StatusOK, resp)
}

func (h *UsersHandler) createProfile(c *gin.Context) {
	userID, ok := h.getUserIDFromClaims(c)
	if !ok {
		return
	}
	var req dto.CreateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}
	prof := &entity.Profiles{IdProfile: helper.Uuid(), IdUser: userID, Name: req.Name, Contact: req.Contact, Address: req.Address, ImageUrl: req.ImageUrl}
	saved, err := h.profile.Create(prof)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to create profile"))
		return
	}
	c.JSON(http.StatusCreated, helper.SuccessResponse("created", saved))
}

func (h *UsersHandler) updateProfile(c *gin.Context) {
	userID, ok := h.getUserIDFromClaims(c)
	if !ok {
		return
	}
	id := c.Param("id")
	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}
	existing, err := h.profile.GetByID(id)
	if err != nil || existing.IdUser != userID {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("profile not found"))
		return
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Contact != nil {
		existing.Contact = *req.Contact
	}
	if req.Address != nil {
		existing.Address = *req.Address
	}
	if req.ImageUrl != nil {
		existing.ImageUrl = *req.ImageUrl
	}
	updated, err := h.profile.Update(id, existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to update profile"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("updated", updated))
}

func (h *UsersHandler) deleteProfile(c *gin.Context) {
	userID, ok := h.getUserIDFromClaims(c)
	if !ok {
		return
	}
	id := c.Param("id")
	existing, err := h.profile.GetByID(id)
	if err != nil || existing.IdUser != userID {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("profile not found"))
		return
	}
	if err := h.profile.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to delete profile"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("deleted", gin.H{"id": id}))
}

func (h *UsersHandler) updateEmail(c *gin.Context) {
	userID, ok := h.getUserIDFromClaims(c)
	if !ok {
		return
	}
	var req dto.UpdateEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}
	user, err := h.Users.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("user not found"))
		return
	}
	user.Email = req.Email
	updated, err := h.Users.Update(userID, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to update email"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("email updated", gin.H{"email": updated.Email}))
}
