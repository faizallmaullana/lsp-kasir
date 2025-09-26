package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/http/dto"
	"faizalmaulana/lsp/http/middleware"
	"faizalmaulana/lsp/http/services"
	"faizalmaulana/lsp/models/entity"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

	// Admin-only user management
	ug := rr.Group("/users")
	ug.POST("", middleware.JWTMiddleware(h.cfg), h.createUserWithProfileAdmin)
	ug.GET("", middleware.JWTMiddleware(h.cfg), h.listUsers)
}

// createUserWithProfileAdmin allows admin to create a new user along with a primary profile
// Body: { email, password, role(optional), profile: { name, contact, address, image_url } }
func (h *UsersHandler) createUserWithProfileAdmin(c *gin.Context) {
	// Auth + role check
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}
	claims, ok := claimsVal.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}
	roleStr, _ := claims["role"].(string)
	if strings.ToLower(roleStr) != "admin" {
		// treat as unauthorized for now (no dedicated Forbidden helper)
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}

	var req dto.CreateUserWithProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}
	setRole := req.Role
	if strings.TrimSpace(setRole) == "" {
		setRole = "cashier"
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to hash password"))
		return
	}

	// Create user
	uid := helper.Uuid()
	u := &entity.Users{
		IdUser:   uid,
		Email:    req.Email,
		Password: string(hashed),
		Role:     setRole,
	}
	createdUser, err := h.Users.Create(u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to create user"))
		return
	}

	// Create profile
	prof := &entity.Profiles{
		IdProfile: helper.Uuid(),
		IdUser:    createdUser.IdUser,
		Name:      req.Profile.Name,
		Contact:   req.Profile.Contact,
		Address:   req.Profile.Address,
		ImageUrl:  req.Profile.ImageUrl,
	}
	createdProfile, err := h.profile.Create(prof)
	if err != nil {
		// best-effort rollback user creation
		_ = h.Users.Delete(createdUser.IdUser)
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to create profile"))
		return
	}

	// Build safe response (omit password)
	resp := gin.H{
		"user": gin.H{
			"id_user":    createdUser.IdUser,
			"email":      createdUser.Email,
			"role":       createdUser.Role,
			"is_deleted": createdUser.IsDeleted,
			"timestamp":  createdUser.Timestamp,
		},
		"profile": createdProfile,
	}
	c.JSON(http.StatusCreated, helper.SuccessResponse("created", resp))
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

// listUsers returns a paginated list of users (admin only)
// Query params: count (default 10, max 100), page (default 1)
func (h *UsersHandler) listUsers(c *gin.Context) {
	// Auth + role check
	claimsVal, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}
	claims, ok := claimsVal.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}
	roleStr, _ := claims["role"].(string)
	if strings.ToLower(roleStr) != "admin" {
		c.JSON(http.StatusUnauthorized, helper.UnauthorizedResponse())
		return
	}

	// Parse pagination
	countQ := c.Query("count")
	pageQ := c.Query("page")
	count, _ := strconv.Atoi(countQ)
	page, _ := strconv.Atoi(pageQ)

	users, err := h.Users.GetAll(count, page)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to list users"))
		return
	}

	// Omit passwords in response
	out := make([]gin.H, 0, len(users))
	for _, u := range users {
		out = append(out, gin.H{
			"id_user":    u.IdUser,
			"email":      u.Email,
			"role":       u.Role,
			"is_deleted": u.IsDeleted,
			"timestamp":  u.Timestamp,
		})
	}

	c.JSON(http.StatusOK, helper.SuccessResponse("OK", out))
}
