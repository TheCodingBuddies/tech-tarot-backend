package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"tech-tarot-backend/cardstack"
)

type CardStackController struct {
	service *cardstack.StackService
}

func NewCardStackController() *CardStackController {
	return &CardStackController{
		service: cardstack.GetService(),
	}
}
func (controller *CardStackController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/image", withCORS(loadImage))
	mux.HandleFunc("/cards", withCORS(draw))
}

func draw(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(cardstack.GetService().Draw3())
}

func loadImage(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		http.Error(writer, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.Atoi(request.URL.Query().Get("id"))
	if err != nil {
		http.Error(writer, "No image url provided or invalid", http.StatusBadRequest)
	}

	url, err := cardstack.GetService().GetImageURL(id)
	if err != nil {
		http.Error(writer, "Id not found", http.StatusNotFound)
		return
	}
	
	image, err := os.Open(url)
	if err != nil {
		http.Error(writer, "Image not found", http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "image/png")
	defer image.Close()

	io.Copy(writer, image)
}
