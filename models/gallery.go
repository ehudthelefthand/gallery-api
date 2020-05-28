package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	Name string
}

type GalleryService interface {
	Create(gallery *Gallery) error
	List() ([]Gallery, error)
	DeleteGallery(id uint) error
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryGorm{
		db: db,
	}
}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) List() ([]Gallery, error) {
	galleries := []Gallery{}
	if err := gg.db.Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (gg *galleryGorm) DeleteGallery(id uint) error {
	return gg.db.Where("id = ?", id).Delete(&Gallery{}).Error
}
