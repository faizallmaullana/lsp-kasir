package handler

import (
	"io"
	"net/http"

	"faizalmaulana/lsp/conf"
	"faizalmaulana/lsp/helper"
	"faizalmaulana/lsp/http/dto"
	"faizalmaulana/lsp/http/middleware"
	"faizalmaulana/lsp/http/services"

	"github.com/gin-gonic/gin"
)

type ImagesHandler struct {
	cfg *conf.Config
	svc services.ImagesService
}

func NewImagesHandler(cfg *conf.Config, svc services.ImagesService) *ImagesHandler {
	return &ImagesHandler{cfg: cfg, svc: svc}
}

func (h *ImagesHandler) Register(rr *gin.RouterGroup) {
	rg := rr.Group("/images")
	rg.POST("/upload", middleware.JWTMiddleware(h.cfg), h.uploadMultipart)
	rg.POST("/upload/base64", middleware.JWTMiddleware(h.cfg), h.uploadBase64)
	rg.GET(":id", h.downloadBlob)
	rg.GET(":id/base64", h.downloadBase64)
	rg.DELETE(":id", middleware.JWTMiddleware(h.cfg), h.delete)
}

func (h *ImagesHandler) uploadMultipart(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse("file is required"))
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse("cannot open file"))
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse("cannot read file"))
		return
	}
	id, stored, err := h.svc.UploadBlob(file.Filename, file.Header.Get("Content-Type"), data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to upload"))
		return
	}
	c.JSON(http.StatusCreated, helper.SuccessResponse("created", dto.UploadResponse{IdImage: id, FileName: stored}))
}

func (h *ImagesHandler) uploadBase64(c *gin.Context) {
	var req dto.UploadBase64Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helper.BadRequestResponse(err.Error()))
		return
	}
	id, stored, err := h.svc.UploadBase64(req.FileName, req.ContentType, req.DataBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to upload"))
		return
	}
	c.JSON(http.StatusCreated, helper.SuccessResponse("created", dto.UploadResponse{IdImage: id, FileName: stored}))
}

func (h *ImagesHandler) downloadBlob(c *gin.Context) {
	id := c.Param("id")
	img, err := h.svc.GetBlob(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("image not found"))
		return
	}
	c.Data(http.StatusOK, img.ContentType, img.Data)
}

func (h *ImagesHandler) downloadBase64(c *gin.Context) {
	id := c.Param("id")
	_, ct, b64, err := h.svc.GetBase64(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.NotFoundResponse("image not found"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("OK", gin.H{"content_type": ct, "data_base64": b64}))
}

func (h *ImagesHandler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, helper.InternalErrorResponse("failed to delete image"))
		return
	}
	c.JSON(http.StatusOK, helper.SuccessResponse("deleted", gin.H{"id": id}))
}
