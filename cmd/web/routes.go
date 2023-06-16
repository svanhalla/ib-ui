package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/svanhalla/ib-ui/internal/models"

	"github.com/gorilla/websocket"

	"github.com/go-chi/chi/v5"
	cors "github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:     []string{"*", "https://*", "http://*"},
		AllowOriginFunc:    nil,
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     nil,
		AllowCredentials:   false,
		MaxAge:             300,
		OptionsPassthrough: false,
		Debug:              false,
	}))

	mux.Use(SessionLoad)

	mux.Get("/", app.Home)
	mux.Get("/ws", Websockets)
	mux.Get("/occasions", app.Occasions)
	mux.Get("/occasions/{uuid}", app.Occasion)
	mux.Delete("/occasions/{uuid}", app.DeleteOccasion)
	mux.Post("/occasions", app.UpdateOccasion)

	mux.Get("/resize-image", app.ResizeForm)
	mux.Post("/api/resize", app.ResizeDo)

	mux.Get("/browse-photos", app.BrowsePhotos)
	mux.Post("/api/browse-photos", app.BrowseDirectory)

	mux.Get("/api/image", app.GetImage)
	mux.Post("/api/generate", app.GenerateOccasion)

	return mux
}

var upgrader = websocket.Upgrader{}

type generateRequest struct {
	UUID    string `json:"uuid"`
	Message string `json:"message"`
	File    string `json:"file"`
}

var messageChannel = make(chan generateRequest)

func Websockets(w http.ResponseWriter, r *http.Request) {

	progresReporter := func(id string, status string) {
		message := generateRequest{
			UUID:    id,
			Message: status,
		}
		messageChannel <- message
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Fel vid uppgradering till WebSocket:", err)
		return
	}
	defer conn.Close()

	// Gorutin för att skicka data från kanalen till WebSocket
	go func() {
		for {
			message := <-messageChannel
			messageBytes, err := json.MarshalIndent(message, "", " ")
			err = conn.WriteMessage(websocket.TextMessage, messageBytes)
			if err != nil {
				log.Println("Fel vid skickande av meddelande till WebSocket:", err)
				break
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Fel vid läsning av meddelande från WebSocket:", err)
			break
		}

		genRequest := generateRequest{}
		err = json.Unmarshal([]byte(message), &genRequest)
		if err != nil {
			log.Println("Failed unmarshal request: ", err)
			break
		}

		answer := generateRequest{}
		answer.UUID = genRequest.UUID

		// check if the file exists

		if _, err := os.Stat(genRequest.File); err == nil {
			answer.Message = "started on server ..."
			// start the process
			// här kan jag fixa implementation
			jsonFile, err := os.Open(genRequest.File)
			if err != nil {
				answer.Message = err.Error()
			}
			// get the definition from file
			definition := models.OccasionDefinition{}

			err = json.NewDecoder(jsonFile).Decode(&definition)
			if err != nil {
				// return fmt.Errorf("failed to decode json file %s: %w", config, err)
				answer.Message = err.Error()
			}

			definition.ProgressReporter = progresReporter
			//definition.Done = make(chan bool)
			//
			definition.WriteProgress("resetted")
			//
			err = definition.GenerateOccasion()
			if err != nil {
				answer.Message = err.Error()
			}
			definition.WriteProgress("DONE")

		} else {
			answer.Message = "File does not exist"
		}
	}
}
