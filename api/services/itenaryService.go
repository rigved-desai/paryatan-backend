package services

import (
	"sort"
	"sync"

	"github.com/rigved-desai/paryatan-backend/api/interfaces"
	"github.com/rigved-desai/paryatan-backend/api/models"
)

// implements ItenaryManager from interfaces pkg
type ItenaryService struct {
	DataAccessor interfaces.ItenaryDataAccessor
}

func (service *ItenaryService) GetItenary(latitude, longitude float64, preferences []string, numberOfDaysAvailable int) (models.Itenary, error) {
	clusterLabel, err := service.DataAccessor.GetClusterByCoordinates(latitude, longitude)

	if err != nil {
		return models.Itenary{}, err
	}

	var placesWithOriginalScores []models.Place
	var minDistance float64
	var maxDistance float64
	var err1, err2 error
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		placesWithOriginalScores, err1 = service.DataAccessor.GetPlacesWithOriginalScores(clusterLabel, latitude, longitude)
	}()

	go func() {
		defer wg.Done()
		minDistance, maxDistance, err2 = service.DataAccessor.GetMinAndMaxDistanceFromUser(clusterLabel, latitude, longitude)
	}()

	wg.Wait()

	if err1 != nil {
		return models.Itenary{}, err1
	}
	if err2 != nil {
		return models.Itenary{}, err2
	}

	placesWithIntermediateScores := service.getPlacesWithIntermediateScores(
		placesWithOriginalScores,
		preferences,
		minDistance,
		maxDistance,
	)

	sort.Slice(
		placesWithIntermediateScores,
		func(i, j int) bool {
			return placesWithIntermediateScores[i].VisitabilityScore > placesWithIntermediateScores[j].VisitabilityScore
		},
	)

	finalScoresAndDistances := service.getPlacesWithFinalScores(placesWithIntermediateScores)

	finalPlacesInItenary := service.choosePlacesAccordingToNumberOfDays(finalScoresAndDistances, numberOfDaysAvailable)

	return models.Itenary{
		Places: finalPlacesInItenary,
	}, nil
}

//below internal functions can be optimized for space by passing objects by reference?
func (service *ItenaryService) getPlacesWithIntermediateScores(originalScoresAndDistances []models.Place, preferences []string, minDistance, maxDistance float64) (placesWithIntermediateScores []models.Place) {
	for _, place := range originalScoresAndDistances {
		personalizationScore := 0.0
		for _, preference := range preferences {
			if preference == place.TypeOfDestination {
				personalizationScore = 1.0
				break
			}
		}
		normalizedDistance := (place.DistanceFromUser - minDistance) / (maxDistance - minDistance)
		newScore := 0.25*place.VisitabilityScore + 0.5*personalizationScore - 0.25*normalizedDistance
		placeWithNewScore := place
		placeWithNewScore.VisitabilityScore = newScore
		placesWithIntermediateScores = append(placesWithIntermediateScores, placeWithNewScore)
	}
	return placesWithIntermediateScores
}

func (service *ItenaryService) getPlacesWithFinalScores(placesWithIntermediateScores []models.Place) (placesWithFinalScores []models.Place) {
	penalties := make(map[string]int)
	for _, place := range placesWithIntermediateScores {
		penaltyFactor := penalties[place.TypeOfDestination]
		penalties[place.TypeOfDestination]++
		placeWithNewScore := place
		// **************** penalty constant (0.25) used here, will probably need to change ****************
		placeWithNewScore.VisitabilityScore = placeWithNewScore.VisitabilityScore - float64(penaltyFactor)*float64(penaltyFactor)*0.25
		placesWithFinalScores = append(placesWithFinalScores, placeWithNewScore)
	}
	return placesWithFinalScores
}

func (service *ItenaryService) choosePlacesAccordingToNumberOfDays(placesWithFinalScores []models.Place, numberOfDaysAvailable int) (placesinItenary []models.Place) {
	sort.Slice(
		placesWithFinalScores,
		func(i, j int) bool {
			return placesWithFinalScores[i].VisitabilityScore > placesWithFinalScores[j].VisitabilityScore
		},
	)
	if numberOfDaysAvailable*2 > len(placesWithFinalScores) {
		return placesWithFinalScores
	}
	return placesWithFinalScores[:2*numberOfDaysAvailable]
}