package models

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	Name      string
	IsPublish bool
	UserID    uint
	Images    []Image
}

type GalleryService interface {
	Create(gallery *Gallery) error
	ListAllPublish() ([]Gallery, error)
	ListByUserID(id uint) ([]Gallery, error)
	GetByID(id uint) (*Gallery, error)
	DeleteGallery(id uint) error
	UpdateGalleryName(id uint, name string) error
	UpdateGalleryPublishing(id uint, isPublish bool) error
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

func (gg *galleryGorm) ListAllPublish() ([]Gallery, error) {
	galleries := []Gallery{}
	err := gg.db.
		Where("is_publish = ?", true).
		Find(&galleries).Error
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(galleries); i++ {
		images := []Image{}
		err := gg.db.
			Where("gallery_id = ?", galleries[i].ID).
			Find(&images).Error
		if err != nil {
			return nil, err
		}
		galleries[i].Images = images
	}
	return galleries, nil
}

func (gg *galleryGorm) ListByUserID(id uint) ([]Gallery, error) {
	galleries := []Gallery{}
	if err := gg.db.
		Where("user_id = ?", id).
		Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

func (gg *galleryGorm) GetByID(id uint) (*Gallery, error) {
	gallery := new(Gallery)
	if err := gg.db.First(gallery, id).Error; err != nil {
		return nil, err
	}
	return gallery, nil
}

func (gg *galleryGorm) DeleteGallery(id uint) error {
	tx := gg.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	err := gg.db.Where("gallery_id = ?", id).Delete(&Image{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	err = gg.db.Where("id = ?", id).Delete(&Gallery{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	idStr := strconv.FormatUint(uint64(id), 10)
	if err := os.RemoveAll(filepath.Join(UploadPath, idStr)); err != nil {
		log.Printf("Fail deleting image files: %v\n", err)
	}

	return tx.Commit().Error
}

func (gg *galleryGorm) UpdateGalleryName(id uint, name string) error {
	return gg.db.Model(&Gallery{}).Where("id = ?", id).Update("name", name).Error
}

func (gg *galleryGorm) UpdateGalleryPublishing(id uint, isPublish bool) error {
	return gg.db.Model(&Gallery{}).Where("id = ?", id).Update("is_publish", isPublish).Error
}
