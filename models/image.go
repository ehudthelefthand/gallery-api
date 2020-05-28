package models

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jinzhu/gorm"
)

const limit = 20
const UploadPath = "upload"

type Image struct {
	gorm.Model
	GalleryID uint   `gorm:"not null"`
	Filename  string `gorm:"not null"`
}

type ImageService interface {
	CreateImages(images []*multipart.FileHeader, galleryID uint) ([]Image, error)
	// Delete(id uint) error
	GetByGalleryID(id uint) ([]Image, error)
}

type imageService struct {
	db *gorm.DB
}

func NewImageService(db *gorm.DB) ImageService {
	return &imageService{db}
}

// Create will insert image to db
func (ims *imageService) CreateImages(files []*multipart.FileHeader, galleryID uint) ([]Image, error) {
	idStr := strconv.FormatUint(uint64(galleryID), 10)
	dir := filepath.Join(UploadPath, idStr)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	tx := ims.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	images := []Image{}
	for _, file := range files {
		dst := filepath.Join(dir, file.Filename)
		if err := saveFile(file, dst); err != nil {
			tx.Rollback()
			return nil, err
		}
		image := Image{
			GalleryID: galleryID,
			Filename:  file.Filename,
		}
		if err := tx.Create(&image).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		images = append(images, image)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return images, nil
}

func saveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// // Delete will delete image from db
// func (ims *ImageService) Delete(image *Image) error {
// 	return ims.DB.Where("id = ?", image.ID).Delete(image).Error
// }

// // GetByID will return image of a given ID
// func (ims *ImageService) GetByID(id uint) (*Image, error) {
// 	image := new(Image)
// 	err := ims.DB.First(image, "id = ?", id).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return image, nil
// }

func (ims *imageService) GetByGalleryID(id uint) ([]Image, error) {
	images := []Image{}
	err := ims.db.
		Order("created_at DESC").
		Where("gallery_id = ?", id).
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}
