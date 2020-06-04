package handlers

import (
	"gallery-api/models"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ImageRes struct {
	ID        uint   `json:"id"`
	GalleryID uint   `json:"gallery_id"`
	Filename  string `json:"filename"`
}

type CreateImageRes struct {
	ImageRes
}

type ImageHandler struct {
	gs  models.GalleryService
	ims models.ImageService
}

func NewImageHandler(gs models.GalleryService, ims models.ImageService) *ImageHandler {
	return &ImageHandler{gs, ims}
}

func (imh *ImageHandler) CreateImage(c *gin.Context) {
	galleryIDStr := c.Param("id")
	galleryID, err := strconv.Atoi(galleryIDStr)
	if err != nil {
		Error(c, 400, err)
		return
	}

	gallery, err := imh.gs.GetByID(uint(galleryID))
	if err != nil {
		Error(c, 400, err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		Error(c, 400, err)
		return
	}

	images, err := imh.ims.CreateImages(form.File["photos"], gallery.ID)
	if err != nil {
		Error(c, 500, err)
		return
	}

	res := []CreateImageRes{}
	for _, img := range images {
		r := CreateImageRes{}
		r.ID = img.ID
		r.GalleryID = gallery.ID
		r.Filename = filepath.Join(models.UploadPath, galleryIDStr, img.Filename)
		res = append(res, r)
	}

	c.JSON(201, res)
}

func (imh *ImageHandler) DeleteImage(c *gin.Context) {
	imageIDStr := c.Param("id")
	id, err := strconv.Atoi(imageIDStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	if err := imh.ims.Delete(uint(id)); err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(http.StatusOK)
}

type ListGalleryImagesRes struct {
	ImageRes
}

func (imh *ImageHandler) ListGalleryImages(c *gin.Context) {
	galleryIDStr := c.Param("id")
	id, err := strconv.Atoi(galleryIDStr)
	if err != nil {
		Error(c, 400, err)
		return
	}

	gallery, err := imh.gs.GetByID(uint(id))
	if err != nil {
		Error(c, 400, err)
		return
	}
	images, err := imh.ims.GetByGalleryID(gallery.ID)
	if err != nil {
		Error(c, http.StatusNotFound, err)
		return
	}
	res := []ListGalleryImagesRes{}
	for _, img := range images {
		r := ListGalleryImagesRes{}
		r.ID = img.ID
		r.GalleryID = gallery.ID
		r.Filename = img.URLPath()
		res = append(res, r)
	}
	c.JSON(http.StatusOK, res)
}
