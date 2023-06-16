package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/google/uuid"

	"github.com/svanhalla/ib-ui/internal/models"

	"github.com/go-chi/chi/v5"
)

// Home displays the home page
func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "home", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) Occasion(w http.ResponseWriter, r *http.Request) {
	theUUID := chi.URLParam(r, "uuid")
	var occasion = &models.OccasionDefinition{}

	if theUUID != "-1" {
		occasions, err := app.Repo.GetOccasions()
		if err != nil {
			app.errorLog.Println("error getting occasions")
			return
		}
		occasion = occasions[theUUID]
	}

	if err := app.renderTemplate(w, r, "edit-occasion", &templateData{
		Data: map[string]interface{}{
			"occasion": occasion,
		},
	}); err != nil {
		app.errorLog.Println(err)
		return
	}
}

func (app *application) DeleteOccasion(w http.ResponseWriter, r *http.Request) {
	theUUID := chi.URLParam(r, "uuid")
	fmt.Println("delete ?", theUUID)
	occasions, err := app.Repo.GetOccasions()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	for _, occasion := range occasions {
		if occasion.UUID == theUUID {
			fmt.Println("found the uuid", occasion.Filename)
			err := os.Remove(occasion.Filename)
			if err != nil {
				fmt.Printf("Kunde inte ta bort filen: %v", err)
				return
			}

			fmt.Printf("Filen %s har tagits bort\n", occasion.Filename)
			return
			break
		}
	}

}

func (app *application) Occasions(w http.ResponseWriter, r *http.Request) {
	occasions, err := app.Repo.GetOccasions()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	occasionList := []*models.OccasionDefinition{}
	for key := range occasions {
		occasionList = append(occasionList, occasions[key])
	}

	sort.Slice(occasionList, func(i, j int) bool {
		return occasionList[i].Name+occasionList[i].UUID > occasionList[j].Name+occasionList[j].UUID
	})

	if err := app.renderTemplate(w, r, "list-occasion", &templateData{
		Data: map[string]interface{}{
			"occasions": occasionList,
		},
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) UpdateOccasion(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	uuidParam := r.Form.Get("uuid")
	if uuidParam != "-1" {
		// new occasion
		uuidParam = uuid.New().String()
	}

	numberOfColumns, _ := strconv.Atoi(r.Form.Get("numberOfColumns"))
	size, _ := strconv.Atoi(r.Form.Get("size"))
	coverSize, _ := strconv.Atoi(r.Form.Get("cover-size"))
	occasion := models.OccasionDefinition{
		UUID:            uuidParam,
		Name:            r.Form.Get("name"),
		Description:     r.Form.Get("description"),
		Root:            r.Form.Get("root"),
		NumberOfColumns: numberOfColumns,
		Title:           r.Form.Get("title"),
		Size:            size,
		Date:            r.Form.Get("date"),
		Cover: models.Part{
			Dir:             r.Form.Get("cover-dir"),
			Name:            r.Form.Get("cover-name"),
			Size:            coverSize,
			NumberOfColumns: 0,
		},
	}

	index := 0
	for true {
		if index > 100 {
			break
		}
		if r.Form.Has(fmt.Sprintf("part-name-%d", index)) {
			size, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("part-size-%d", index)))
			noc, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("part-cols-%d", index)))
			part := models.Part{
				Dir:             r.Form.Get(fmt.Sprintf("part-dir-%d", index)),
				Name:            r.Form.Get(fmt.Sprintf("part-name-%d", index)),
				Size:            size,
				NumberOfColumns: noc,
			}
			occasion.Parts = append(occasion.Parts, part)
		}
		index++
	}

	err = app.Repo.SaveOccasion(&occasion)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/occasions", http.StatusSeeOther)
}

func (app *application) ResizeForm(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "perform-resize", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) BrowsePhotos(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "browse-photos", &templateData{}); err != nil {
		app.errorLog.Println(err)
	}
}
