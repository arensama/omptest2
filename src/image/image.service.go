package image

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/arensama/testapi2/src/db"
	"github.com/arensama/testapi2/src/errWithStatus"
	"github.com/arensama/testapi2/src/model"
)

type ImageService struct {
	db *db.DB
}

func ServiceInit(db *db.DB) *ImageService {
	db.Migrate(model.Image{})
	db.Migrate(model.Chunk{})
	return &ImageService{
		db: db,
	}
}

func (s *ImageService) CreateImage(image model.Image) (model.Image, error) {
	db := s.db.Db
	var exists model.Image
	err := db.Where("Uid = ?", image.Uid).First(&exists)
	if err.Error == nil {
		return model.Image{},
			errWithStatus.StatusErr{
				Status:  http.StatusConflict,
				Message: fmt.Sprintf("image with sha256 %v already exists!", image.Uid),
			}
	}
	err = db.Create(&image)
	if err.Error != nil {
		return model.Image{}, err.Error
	}
	return image, nil
}

func (s *ImageService) UploadChunk(uid string, chunk model.Chunk) (model.Chunk, error) {
	db := s.db.Db

	// check if image is created
	var image model.Image
	err := db.Where("Uid = ?", uid).First(&image)
	if err.Error != nil {
		return model.Chunk{},
			errWithStatus.StatusErr{
				Status:  http.StatusNotFound,
				Message: fmt.Sprintf("image with sha256 %v already exists!", uid),
			}
	}

	// check if uploaded chunk already exists
	// var exists model.Chunk
	// if err := db.Where("orde = ? ", chunk.Order).First(&exists).Error; err != nil {
	// 	if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 		fmt.Println("err", err)
	// 		return model.Chunk{},
	// 			errWithStatus.StatusErr{
	// 				Status:  http.StatusInternalServerError,
	// 				Message: "internal server error!",
	// 			}
	// 	}
	// } else {
	// 	return model.Chunk{},
	// 		errWithStatus.StatusErr{
	// 			Status:  http.StatusConflict,
	// 			Message: fmt.Sprintf("chunk with id %v already exists!", chunk.Order),
	// 		}
	// }
	chunk.ImageId = image.ID
	chunk.Image = image
	err = db.Create(&chunk)
	if err.Error != nil {
		return model.Chunk{}, errWithStatus.StatusErr{
			Status:  http.StatusInternalServerError,
			Message: err.Error.Error(),
		}
	}
	return chunk, nil
}

func (s *ImageService) DownloadImage(uid string) (string, error) {
	db := s.db.Db
	var image model.Image
	err := db.Preload("Chunks").Where("uid = ?", uid).First(&image)
	if err.Error != nil {
		return "", err.Error
	}
	chunks := image.Chunks
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Order < chunks[j].Order
	})
	wholeImage := ""
	wholeImageSize := 0
	for i := 0; i < len(chunks); i++ {
		chunk := chunks[i]
		wholeImage += chunk.Data
		wholeImageSize += chunk.Size
	}
	return wholeImage, nil
}
