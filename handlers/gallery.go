package handlers

import (
	"gallery-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GalleryHandler struct {
	gs models.GalleryService
}

func NewGalleryHandler(gs models.GalleryService) *GalleryHandler {
	return &GalleryHandler{gs}
}

type CreateGallery struct {
	Name string
}

func (gh *GalleryHandler) Create(c *gin.Context) {
	data := new(CreateGallery)
	if err := c.BindJSON(data); err != nil {
		Error(c, 400, err)
		return
	}
	gallery := new(models.Gallery)
	gallery.Name = data.Name
	if err := gh.gs.Create(gallery); err != nil {
		Error(c, 500, err)
		return
	}
	c.JSON(201, gin.H{
		"id":   gallery.ID,
		"name": gallery.Name,
	})
}

type Gallery struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

func (gh *GalleryHandler) List(c *gin.Context) {
	data, err := gh.gs.List()
	if err != nil {
		Error(c, 500, err)
		return
	}
	galleries := []Gallery{}
	for _, a := range data {
		galleries = append(galleries, Gallery{
			ID:   a.ID,
			Name: a.Name,
		})
	}
	c.JSON(200, galleries)
}

func (gh *GalleryHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	err = gh.gs.DeleteGallery(uint(id))
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}

type UpdateNameReq struct {
	Name string
}

func (gh *GalleryHandler) UpdateName(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	req := new(UpdateNameReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	err = gh.gs.UpdateGalleryName(uint(id), req.Name)
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}

type UpdateStatusReq struct {
	IsPublish bool
}

func (gh *GalleryHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		Error(c, 400, err)
		return
	}
	req := new(UpdateStatusReq)
	if err := c.BindJSON(req); err != nil {
		Error(c, 400, err)
		return
	}
	err = gh.gs.UpdateGalleryStatus(uint(id), req.IsPublish)
	if err != nil {
		Error(c, 500, err)
		return
	}
	c.Status(204)
}
