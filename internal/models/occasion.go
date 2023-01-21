package models

import (
	"encoding/json"
	"os"
)

type OccasionDefinition struct {
	UUID            string `json:"uuid"`
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
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
