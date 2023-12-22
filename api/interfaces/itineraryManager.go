package interfaces

import "github.com/rigved-desai/paryatan-backend/api/models"

// will be imnplemented by ItineraryService
type ItineraryManager interface {
	GetItinerary(startLocationName, startLocationCity string, latitude, longitude float64, preferences []string, numberOfDaysAvailable int) (models.Itinerary, error)
}
