package picture

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/h2non/filetype"
)

type Picture struct {
	Name       string
	Path       string
	GroupName  string
	ScaledPath string
	Type       string

	// add date time
}

func (p *Picture) Resize(imageSize int) error {
	isDir, err := isDirectory(p.Path)
	if err != nil {
		return fmt.Errorf("failed to check if %q is a directory: %w", p.Path, err)
	}

	if isDir {
		dirEntries, err := os.ReadDir(p.Path)
		if err != nil {
			return fmt.Errorf("failed to read dir %q: %w", p.Path, err)
		}

		for _, theFile := range dirEntries {
			// ignore hidden files and directories
			if strings.HasPrefix(theFile.Name(), ".") || theFile.IsDir() {
				continue
			}

			imageFile := filepath.Join(p.Path, theFile.Name())

			err = scaleImageFile(imageFile, imageSize)
			if err != nil {
				return fmt.Errorf("failed to scale %q: %w", imageFile, err)
			}
		}

		return nil
	}

	err = scaleImageFile(p.Path, imageSize)
	if err != nil {
		return fmt.Errorf("failed to scale %q: %w", p.Path, err)
	}

	return nil
}

// isDirectory determines if a file represented.
func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to check if %q is directory: %w", path, err)
	}

	return fileInfo.IsDir(), nil
}

func scaleImageFile(imageFullpath string, imageSize int) error {
	imageType := getFileType(imageFullpath)
	if imageType != imageStr && imageType != videoStr {
		return nil
	}

	// resize the image
	if imageType == imageStr {
		err := scaleImage(imageFullpath, imageSize)
		if err != nil {
			return fmt.Errorf("failed to scale image %q: %w", imageFullpath, err)
		}
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

func scaleImage(filename string, imageSize int) error {
	imageBytes, err := getImageBytes(filename)
	if err != nil {
		return fmt.Errorf("failed to get bytes for image %q: %w", filename, err)
	}

	imageConfig, _, err := image.DecodeConfig(bytes.NewBuffer(imageBytes))
	if err != nil {
		return fmt.Errorf("failed to create image config %q: %w", filename, err)
	}

	// landscape
	if imageConfig.Width > imageConfig.Height {
		err = resize(imageSize, 0, filename)
		if err != nil {
			return fmt.Errorf("failed to resize landscape image %q: %w", filename, err)
		}

		return nil
	}

	// portrait
	_ = resize(0, imageSize, filename)

	return nil
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
