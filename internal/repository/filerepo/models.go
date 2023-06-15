package filerepo

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gosimple/slug"
	"github.com/svanhalla/ib-ui/internal/models"
)

// DBModel is the type for database connection values
type FileModel struct {
	Dir string
}

func (m *FileModel) RootDir() string {
	if strings.HasPrefix(m.Dir, "~") {
		userHome, _ := os.UserHomeDir()
		m.Dir = strings.ReplaceAll(m.Dir, "~", userHome)
	}
	return m.Dir
}

func (m *FileModel) GetOccasions() ([]*models.OccasionDefinition, error) {
	files, err := os.ReadDir(m.RootDir())
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %q: %w", m.Dir, err)
	}

	occasions := []*models.OccasionDefinition{}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		occasion, err := models.NewOccasionFromFile(path.Join(m.Dir, f.Name()))
		if err != nil {
			return nil, err
		}
		occasion.Filename = path.Join(m.Dir, f.Name())

		occasions = append(occasions, occasion)
	}
	return occasions, nil
}

func (m *FileModel) SaveOccasion(occasion *models.OccasionDefinition) error {
	marshalIndent, err := json.MarshalIndent(occasion, "", " ")
	if err != nil {
		return err
	}

	// create a file name
	fileName := path.Join(m.RootDir(), fmt.Sprintf("%s.json", slug.Make(occasion.Name)))

	err = os.WriteFile(fileName, marshalIndent, 0644)
	if err != nil {
		return err
	}
	return nil
}
