package image

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arensama/testapi2/src/errWithStatus"
	"github.com/arensama/testapi2/src/model"
	"github.com/gorilla/mux"
)

type ImageController struct {
	router       *mux.Router
	imageService *ImageService
}

func Init(imageService *ImageService) *ImageController {
	c := ImageController{
		router:       mux.NewRouter(),
		imageService: imageService,
	}
	c.router.HandleFunc("/image", c.createImage).Methods("POST")
	c.router.HandleFunc("/image/{uid}/chunks", c.uploadChunk).Methods("POST")
	c.router.HandleFunc("/image/{uid}", c.downloadImage).Methods("GET")
	return &c
}

func (c *ImageController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.router.ServeHTTP(w, r)
}
func (c *ImageController) createImage(w http.ResponseWriter, r *http.Request) {
	var image model.Image
	if err := json.NewDecoder(r.Body).Decode(&image); err != nil {
		http.Error(w, "Malformed request", http.StatusBadRequest)
		return
	}
	createdImage, err := c.imageService.CreateImage(image)
	if err != nil {
		http.Error(w, err.Error(), err.(errWithStatus.StatusErr).Status)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdImage)
}

func (c *ImageController) uploadChunk(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]
	var chunk model.Chunk
	if err := json.NewDecoder(r.Body).Decode(&chunk); err != nil {
		http.Error(w, "Malformed request", http.StatusBadRequest)
		return
	}
	fmt.Println("decoded", chunk.Order)
	createdChunk, err := c.imageService.UploadChunk(uid, chunk)
	if err != nil {
		http.Error(w, err.Error(), err.(errWithStatus.StatusErr).Status)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdChunk)
}
func (c *ImageController) downloadImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uid := vars["uid"]
	image, err := c.imageService.DownloadImage(uid)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	w.Write([]byte(image))
}
