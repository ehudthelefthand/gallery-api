package models

import (
	"io"
	"log"
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

func (img *Image) FilePath() string {
	idStr := strconv.FormatUint(uint64(img.GalleryID), 10)
	return filepath.Join(UploadPath, idStr, img.Filename)
}

type ImageService interface {
	CreateImages(images []*multipart.FileHeader, galleryID uint) ([]Image, error)
	Delete(id uint) error
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
		image := Image{
			GalleryID: galleryID,
			Filename:  file.Filename,
		}
		if err := tx.Create(&image).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		images = append(images, image)

		if err := saveFile(file, image.FilePath()); err != nil {
			tx.Rollback()
			return nil, err
		}
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

func (ims *imageService) Delete(id uint) error {
	image, err := ims.GetByID(id)
	if err != nil {
		return err
	}
	err = os.Remove(image.FilePath())
	if err != nil {
		log.Printf("Fail deleting image: %v\n", err)
	}
	return ims.db.Where("id = ?", id).Delete(&Image{}).Error
}

// GetByID will return image of a given ID
func (ims *imageService) GetByID(id uint) (*Image, error) {
	image := new(Image)
	err := ims.db.First(image, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return image, nil
}

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
