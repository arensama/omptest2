package model

import "gorm.io/gorm"

type Image struct {
	gorm.Model
	Uid       string  `json:"sha256" gorm:"unique_index;not null"`
	Size      int     `json:"size"`
	ChunkSize int     `json:"chunk_size"`
	Chunks    []Chunk `json:"chunks" gorm:"foreignKey:ImageId"`
}
type Chunk struct {
	gorm.Model
	ID      uint   `json:"-" gorm:"column:id"`
	Order   int    `json:"id" gorm:"column:order"`
	Size    int    `json:"size" gorm:"column:size"`
	Data    string `json:"data" gorm:"column:data"`
	ImageId uint   `json:"-"`
	Image   Image  `json:"-"`
}
