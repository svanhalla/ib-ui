package models

import (
	"os"
	"strings"
)

type ImageDirectory struct {
	Directory string   `json:"directory"`
	Files     []string `json:"files"`
	Dirs      []string `json:"dirs"`
}

func NewImageDir(dir string) (*ImageDirectory, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	imageDir := ImageDirectory{
		Directory: dir,
	}

	for _, f := range files {
		// skip hidden files/dirs
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		if f.IsDir() {
			if err != nil {
				return nil, err
			}
			imageDir.Dirs = append(imageDir.Dirs, f.Name())
			continue
		}
		imageDir.Files = append(imageDir.Files, f.Name())
	}

	return &imageDir, nil
}
