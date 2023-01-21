package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/disintegration/imaging"

	"github.com/svanhalla/ib-ui/internal/models"

	"github.com/svanhalla/ib-ui/internal/picture"
)

func (app *application) ResizeDo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	size, err := strconv.Atoi(r.Form.Get("size"))

	pic := picture.Picture{
		Path: r.Form.Get("image-path"),
	}

	err = pic.Resize(size)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	if err := app.renderTemplate(w, r, "perform-resize", &templateData{
		Data: map[string]interface{}{
			"picture": pic,
		},
	}); err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) BrowseDirectory(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	err := app.ReadJSON(w, r, &data)
	if err != nil {
		_ = app.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	directoryToScan := data["directory"]

	d, err := models.NewImageDir(directoryToScan)
	if err != nil {
		_ = app.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	response := JSONResponse{
		Error:   false,
		Message: "tjo",
		Data:    d,
	}

	err = app.WriteJSON(w, http.StatusOK, response)
	if err != nil {
		fmt.Println("ERROR ", err)
	}
}

func (app *application) GetImage(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	scaleString := r.URL.Query().Get("scale")

	imageBuf, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if scaleString == "" {
		w.Header().Set("Content-Type", "image/jpg")
		_, _ = w.Write(imageBuf)
		return
	}

	scale, _ := strconv.Atoi(scaleString)

	fmt.Println(path)
	fmt.Println(scale)

	theImage, _, err := image.Decode(bytes.NewBuffer(imageBuf))
	if err != nil {
		app.errorLog.Println(fmt.Errorf("failed to decode bytes for image %s ;%w", path, err))
		_ = app.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	imageConfig, _, err := image.DecodeConfig(bytes.NewBuffer(imageBuf))
	if err != nil {
		app.errorLog.Println(fmt.Errorf("failed to decode config for image %s ;%w", path, err))
		log.Fatal(err)
	}

	var resizedImage *image.NRGBA

	// landscape
	if imageConfig.Width > imageConfig.Height {
		fmt.Println("landscape")
		resizedImage = imaging.Resize(theImage, scale, 0, imaging.Lanczos)
	} else {
		fmt.Println("portrait")
		resizedImage = imaging.Resize(theImage, 0, scale, imaging.Lanczos)
	}

	opt := jpeg.Options{
		Quality: 80,
	}

	w.Header().Set("Content-Type", "image/jpg")
	err = jpeg.Encode(w, resizedImage, &opt)
	if err != nil {
		app.errorLog.Println(fmt.Errorf("failed to encode image %s ;%w", path, err))
		log.Fatal(err)
	}

	// _, _ = w.Write(imageBuf)

}
