package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/google/uuid"

	"github.com/svanhalla/ib-ui/internal/models"

	"github.com/svanhalla/ib-ui/internal/picture"
	"github.com/urfave/cli/v2"
)

func main() {
	if err := Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func Run() error {
	app := CreateApp()

	err := app.Run(os.Args)
	if err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}

func CreateApp() cli.App {
	return cli.App{
		Commands: []*cli.Command{
			{
				Name:  "heic-to-jpg",
				Usage: "converts heic file to jpg",
				// Action: ConvertHeic,
			},
			{
				Name:  "image-date",
				Usage: "shows the exif date for image",
				// Action: ImageDate,
			},
			{
				Name:  "get-dates",
				Usage: "scans the directory path and list images date",
				// Action: GetDates,
			},
			{
				Name:  "list-images",
				Usage: "scans the directory path and list file name in date order",
				// Action: ListImages,
			},
			{
				Name:  "scan",
				Usage: "scans the directory url",
				// Action: Scan,
			},
			{
				Name:  "download",
				Usage: "downloads image using url, file name will be base of url",
				// Action: GetImage,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "path",
						Usage: "path with full path where to put downloaded file",
						Value: "/tmp",
					},
				},
			},
			{
				Name:   "resize",
				Usage:  "resizes the image",
				Action: ResizeImage,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "size",
						Usage: "the size for the scaled image",
						Value: 1024,
					},
				},
			},
			{
				Name:  "server",
				Usage: "starts server on localhost",
				// Action: Server,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "port",
						Usage: "port for server to listen on",
						Value: 1964,
					},
				},
			},
			{
				Name:   "generate-occasion",
				Usage:  "generates a occasion page from images in a directory",
				Action: GenerateOccasion,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "config",
						Usage: "path to config file to use",
					},
				},
			},
			{
				Name:  "generate-structure",
				Usage: "generates a directory structure",
				// Action: GenerateStructure,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Usage:    "path to config file to use",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "dir",
						Usage:    "directory where to put the structure",
						Required: true,
					},
				},
			},
			{
				Name:   "generate-example",
				Usage:  "generate example json",
				Action: ExampleJSON,
			},
		},
	}
}

func ResizeImage(cliCtx *cli.Context) error {
	imageSize := cliCtx.Int("size")

	for _, fullpath := range cliCtx.Args().Slice() {
		imageToResize := picture.Picture{
			Name:       "",
			Path:       fullpath,
			GroupName:  "",
			ScaledPath: "",
			Type:       "",
		}

		err := imageToResize.Resize(imageSize)
		if err != nil {
			return fmt.Errorf("failed to resize picture: %w", err)
		}
	}
	return nil
}

func ExampleJSON(_ *cli.Context) error {
	theUUID, _ := uuid.NewUUID()
	var Definition = models.OccasionDefinition{
		UUID:            theUUID.String(),
		Name:            "generated example",
		Description:     "generated example",
		Root:            "/tmp/images/",
		NumberOfColumns: 4,
		Title:           "the title",
		Size:            256,
		Date:            "the date",
		Location:        "the location ",
		Cover: models.Part{
			Dir:  "0-cover",
			Size: 1024,
		},
		Parts: []models.Part{
			{
				Dir:  "1-preparations",
				Name: "Förberedelse",
			},
			{
				Dir:  "2-cermony",
				Name: "Cermoni",
			},
			{
				Dir:  "3-mingle",
				Name: "Mingel",
			},
			{
				Dir:  "4-dinner",
				Name: "Middag",
			},
			{
				Dir:  "5-party",
				Name: "Fest",
			},
			{
				Dir:  "6-portrait",
				Name: "Porträtt",
			},
		},
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home dir: %w", err)
	}

	ibDir := path.Join(homeDir, ".ib")
	if _, err := os.Stat(ibDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(ibDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to crete dir %s:%w", ibDir, err)
		}
	}

	marshalIndent, err := json.MarshalIndent(Definition, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal to json:%w", err)
	}

	fileToWrite := path.Join(ibDir, "example.json")
	err = os.WriteFile(fileToWrite, marshalIndent, 0600)
	if err != nil {
		return fmt.Errorf("failed to write file %s:%w", fileToWrite, err)
	}

	fmt.Printf("created file '%s'\n", fileToWrite)
	return nil
}

func GenerateOccasion(cliCtx *cli.Context) error {
	config := cliCtx.String("config")

	jsonFile, err := os.Open(config)
	if err != nil {
		return fmt.Errorf("failed to open file %s:%w", config, err)
	}

	// get the definition from file
	definition := models.OccasionDefinition{}

	err = json.NewDecoder(jsonFile).Decode(&definition)
	if err != nil {
		return fmt.Errorf("failed to decode json file %s: %w", config, err)
	}

	err = definition.GenerateOccasion()
	if err != nil {
		return fmt.Errorf("failed to generate occasion: %w", err)
	}

	return nil
}
