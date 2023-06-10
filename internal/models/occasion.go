package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/disintegration/imaging"
	"github.com/gosimple/slug"
	"github.com/h2non/filetype"
	"github.com/svanhalla/ib-ui/static"
)

type OccasionDefinition struct {
	UUID            string `json:"uuid"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Root            string `json:"root,omitempty"` // the root directory
	NumberOfColumns int    `json:"numberOfColumns"`
	Title           string `json:"title,omitempty"`    // the title on the cover image
	Size            int    `json:"size,omitempty"`     // the size for scales images
	Date            string `json:"date,omitempty"`     // date for the occasion
	Location        string `json:"location,omitempty"` // location for the occasion
	Cover           Part   `json:"cover"`              // the cover image(s)
	Parts           []Part `json:"parts,omitempty"`    // the page parts
}

type Part struct {
	Dir             string `json:"dir,omitempty"`
	Name            string `json:"name,omitempty"`
	Size            int    `json:"size,omitempty"`
	NumberOfColumns int    `json:"numberOfColumns"`
}

func NewOccasionFromFile(file string) (*OccasionDefinition, error) {
	fileReader, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = fileReader.Close()
	}()

	occasion := &OccasionDefinition{}
	err = json.NewDecoder(fileReader).Decode(occasion)
	if err != nil {
		return nil, err
	}
	return occasion, nil
}

type Image struct {
	Name       string
	Path       string
	GroupName  string
	ScaledPath string
	Type       string

	// add date time
}

func (i *Image) GetImageBytes(p string) ([]byte, error) {
	return getImageBytes(p + i.Path)
}

func getImageBytes(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s:%w", filename, err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s:%w", filename, err)
	}

	return data, nil
}

func (i *Image) GetImageConfig(p string) (image.Config, error) {
	imageBytes, err := i.GetImageBytes(p)
	if err != nil {
		return image.Config{}, fmt.Errorf("failed to get bytes for image: %w", err)
	}

	imageConfig, _, err := image.DecodeConfig(bytes.NewBuffer(imageBytes))
	if err != nil {
		return imageConfig, fmt.Errorf("failed to create image config: %w", err)
	}

	return imageConfig, nil
}

func (i *Image) GetImage(p string) (image.Image, error) {
	imageBytes, err := i.GetImageBytes(p)
	if err != nil {
		return nil, fmt.Errorf("failed to get image bytes: %w", err)
	}

	theImage, _, err := image.Decode(bytes.NewBuffer(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return theImage, nil
}

func (i *Image) Resize(width int, height int, p string) (string, error) {
	extension := filepath.Ext(i.Path)
	newFilename := i.Name[:len(i.Name)-len(filepath.Ext(i.Name))]
	ScaledFile := fmt.Sprintf("%s/scaled/%s-%dx%d%s", i.GroupName, newFilename, width, height, extension)

	// check if exists
	if _, err := os.Stat(p + ScaledFile); err == nil {
		return ScaledFile, nil
	}

	// create dir if not exists
	if _, err := os.Stat(p + i.GroupName + "/scaled/"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(p+i.GroupName+"/scaled/", os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create dir %s: %w", p+i.GroupName+"/scaled/", err)
		}
	}

	theImage, err := i.GetImage(p)
	if err != nil {
		return "", fmt.Errorf("failed to get image: %w", err)
	}

	resizedImage := imaging.Resize(theImage, width, height, imaging.Lanczos)

	resizedImageFile, err := os.Create(p + ScaledFile)
	if err != nil {
		return "", fmt.Errorf("failed to create file :%w", err)
	}

	opt := jpeg.Options{
		Quality: 80,
	}

	err = jpeg.Encode(resizedImageFile, resizedImage, &opt)
	if err != nil {
		return ScaledFile, fmt.Errorf("failed to decode: %w", err)
	}

	return ScaledFile, nil
}

func resize(width int, height int, filepath string) error {
	filename := path.Base(filepath)
	extension := path.Ext(filepath)
	newFilename := filename[:len(filename)-len(path.Ext(filename))]
	dir := path.Dir(filepath)
	scaledFilepath := fmt.Sprintf("%s/scaled/%s-%dx%d%s", dir, newFilename, width, height, extension)

	// check if exists
	if _, err := os.Stat(scaledFilepath); err == nil {
		return nil
	}

	// create dir if not exists
	if _, err := os.Stat(path.Join(dir, "scaled")); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path.Join(dir, "scaled"), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create dir %s: %w", path.Join(dir, "scaled"), err)
		}
	}

	imageBytes, err := getImageBytes(filepath)
	if err != nil {
		return fmt.Errorf("failed to get image bytes: %w", err)
	}

	theImage, _, err := image.Decode(bytes.NewBuffer(imageBytes))
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	resizedImage := imaging.Resize(theImage, width, height, imaging.Lanczos)

	resizedImageFile, err := os.Create(scaledFilepath)
	if err != nil {
		return fmt.Errorf("failed to create file :%w", err)
	}

	opt := jpeg.Options{
		Quality: 80,
	}

	err = jpeg.Encode(resizedImageFile, resizedImage, &opt)
	if err != nil {
		return fmt.Errorf("failed to decode: %w", err)
	}

	return nil
}

type WeddingPart struct {
	Name   string
	ID     string
	Images map[int][]Image
}

type Occasion struct {
	Title          string
	SlideShow      bool
	NumberOfSlides int
	Cover          map[int][]Image
	Date           string
	Location       string
	Parts          []WeddingPart
	Columns        int
}

func createImageMap(imageSize int, numberOfColumns int, path string, directoryName string) map[int][]Image {
	wantedMap := make(map[int][]Image, numberOfColumns)

	for i, theImage := range getImagesInDir(path + directoryName) {
		if theImage.Name == ".DS_Store" {
			continue
		}

		// resize the image
		if theImage.Type == imageStr {
			imageConfig, err := theImage.GetImageConfig(path)
			if err != nil {
				panic("failed to get image config")
			}

			// landscape
			if imageConfig.Width > imageConfig.Height {
				theImage.ScaledPath, _ = theImage.Resize(imageSize, 0, path)
				wantedMap[i%numberOfColumns] = append(wantedMap[i%numberOfColumns], theImage)

				continue
			}

			// portrait
			theImage.ScaledPath, _ = theImage.Resize(0, imageSize, path)
			wantedMap[i%numberOfColumns] = append(wantedMap[i%numberOfColumns], theImage)
		}

		if theImage.Type == videoStr {
			theImage.ScaledPath = theImage.Path
			wantedMap[i%numberOfColumns] = append(wantedMap[i%numberOfColumns], theImage)
		}
	}

	return wantedMap
}

func getImagesInDir(thePath string) []Image {
	var images []Image

	files, err := os.ReadDir(thePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.Name() == ".DS_Store" || file.IsDir() {
			continue
		}

		imageType := getFileType(filepath.Join(thePath, file.Name()))
		if imageType != imageStr && imageType != videoStr {
			continue
		}

		images = append(images, Image{
			Name:      file.Name(),
			Path:      filepath.Join(path.Base(thePath), file.Name()),
			GroupName: path.Base(thePath),
			Type:      imageType,
		})
	}

	return images
}

func (d OccasionDefinition) GenerateOccasion() error {
	occasion := Occasion{
		Columns:  100 / d.NumberOfColumns,
		Title:    d.Title,
		Date:     d.Date,
		Location: d.Location,
		Cover:    createImageMap(d.Cover.Size, 1, d.Root, d.Cover.Dir),
		Parts:    make([]WeddingPart, 0),
	}

	occasion.SlideShow = len(occasion.Cover[0]) > 1
	occasion.NumberOfSlides = len(occasion.Cover[0])

	for _, part := range d.Parts {
		size := d.Size
		if part.Size > 0 {
			size = part.Size
		}

		occasion.Parts = append(occasion.Parts, WeddingPart{
			Name:   part.Name,
			ID:     slug.Make(part.Name),
			Images: createImageMap(size, d.NumberOfColumns, d.Root, part.Dir),
		})
	}

	funcMap := template.FuncMap{
		// The name "inc" is what the function will be called in the template text.
		"inc": func(i int) int {
			return i + 1
		},
	}

	tmpl := template.New("occasion").Funcs(funcMap)
	tmpl = template.Must(tmpl.ParseFS(static.GetFS(), []string{"occasion-template.html"}...))

	// tmpl := template.Must(template.ParseFS(static.GetFS(), []string{"occasion-template.html"}...))

	myfile, err := os.Create(d.Root + "index.html")
	if err != nil {
		return fmt.Errorf("failed to create file %w", err)
	}

	defer func() {
		_ = myfile.Close()
	}()

	err = tmpl.Execute(myfile, occasion)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}

func getFileType(filename string) string {
	file, _ := os.Open(filename)

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)

	_, err := file.Read(head)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read file %s: %s", filename, err.Error())

		return ""
	}

	if filetype.IsImage(head) {
		return imageStr
	}

	if filetype.IsVideo(head) {
		return videoStr
	}

	return ""
}

const imageStr string = "image"
const videoStr string = "video"
