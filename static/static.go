package static

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
)

//go:embed *
var Assets embed.FS

func GetFile(name string) ([]byte, error) {
	if _, err := os.Stat("static/"); os.IsNotExist(err) {
		theBytes, err := Assets.ReadFile(name)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", name, err)
		}

		return theBytes, nil
	}

	theBytes, err := os.ReadFile("static/" + name)
	if err != nil {
		return nil, fmt.Errorf("failed to read from static/%s: %w", name, err)
	}

	return theBytes, nil
}

func GetFS() fs.FS {
	if _, err := os.Stat("static/"); os.IsNotExist(err) {
		return Assets
	}

	return os.DirFS("./static/")
}
