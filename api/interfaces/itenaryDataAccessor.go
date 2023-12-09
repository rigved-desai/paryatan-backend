package interfaces

import "github.com/rigved-desai/paryatan-backend/api/models"

// will be implemented by the data acess layer present in the datastore pkg
type ItenaryDataAccessor interface {
	GetClusterByCoordinates(latitude, longitude float64) (int, error)
	GetPlacesWithOriginalScores(cluster int, latitude, longitude float64) ([]models.Place, error)
	GetMinAndMaxDistanceFromUser(cluster int, latitude, longitude float64) (float64, float64, error)
}