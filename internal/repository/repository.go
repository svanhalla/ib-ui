package repository

import "github.com/svanhalla/ib-ui/internal/models"

type Repo interface {
	GetOccasions() (map[string]*models.OccasionDefinition, error)
	SaveOccasion(occasion *models.OccasionDefinition) error
}
