package main

import (
	"fmt"
	"os"

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
				Name:  "generate-occasion",
				Usage: "generates a occasion page from images in a directory",
				// Action: GenerateOccasion,
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
				Name:  "generate-example",
				Usage: "generate example json",
				// Action: ExampleJSON,
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
