package models

import (
	"gallery-api/rand"
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

func (ims *imageService) CreateImages(files []*multipart.FileHeader, galleryID uint) ([]Image, error) {
	idStr := strconv.FormatUint(uint64(galleryID), 10)
	dir := filepath.Join(UploadPath, idStr)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("create gallery dir error: %v\n", err)
		return nil, err
	}

	tx := ims.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			log.Println("rollback due to error while uploading photo")
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Printf("transaction error: %v\n", err)
		return nil, err
	}

	images := []Image{}
	for _, file := range files {
		generated, err := rand.GetToken()
		if err != nil {
			log.Printf("error generating filename: %v\n", err)
			tx.Rollback()
			return nil, err
		}
		ext := filepath.Ext(file.Filename)
		image := Image{
			GalleryID: galleryID,
			Filename:  generated[:len(generated)-1] + ext,
		}
		if err := tx.Create(&image).Error; err != nil {
			log.Printf("create image error: %v\n", err)
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
		log.Printf("transaction commit: %v\n", err)
		return nil, err
	}

	return images, nil
}

func saveFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		log.Printf("saveFile - open file error: %v\n", err)
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		log.Printf("saveFile - create destination file error: %v\n", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		log.Printf("saveFile - copy file error: %v\n", err)
	}
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
