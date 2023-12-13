package interfaces

import "github.com/rigved-desai/paryatan-backend/api/models"

// will be imnplemented by ItenaryService
type ItenaryManager interface {
	GetItenary(startLocationName, startLocationCity string,  latitude, longitude float64, preferences []string, numberOfDaysAvailable int) (models.Itenary, error)
}