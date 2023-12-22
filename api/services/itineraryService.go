package services

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/rigved-desai/paryatan-backend/api/interfaces"
	"github.com/rigved-desai/paryatan-backend/api/models"
)

// implements ItineraryManager from interfaces pkg
type ItineraryService struct {
	DataAccessor interfaces.ItineraryDataAccessor
}

func (service *ItineraryService) GetItinerary(startLocationName, startLocationCity string, latitude, longitude float64, preferences []string, numberOfDaysAvailable int) (models.Itinerary, error) {
	clusterLabel, err := service.DataAccessor.GetClusterByCoordinates(latitude, longitude)

	if err != nil {
		return models.Itinerary{}, err
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
		return models.Itinerary{}, err1
	}
	if err2 != nil {
		return models.Itinerary{}, err2
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

	finalPlacesInItinerary := service.choosePlacesAccordingToNumberOfDays(finalScoresAndDistances, numberOfDaysAvailable)

	finalDayPlans := service.setDayPlans(finalPlacesInItinerary, numberOfDaysAvailable, startLocationName, startLocationCity, latitude, longitude)

	return models.Itinerary{
		DayPlans: finalDayPlans,
	}, nil
}

// below internal functions can be optimized for space by passing objects by reference?
func (service *ItineraryService) getPlacesWithIntermediateScores(originalScoresAndDistances []models.Place, preferences []string, minDistance, maxDistance float64) (placesWithIntermediateScores []models.Place) {
	VISIBILITY_SCORE_WEIGHT, _ := strconv.ParseFloat(os.Getenv("VISIBILITY_SCORE_WEIGHT"), 64)
	PERSONALIZATION_SCORE_WEIGHT, _ := strconv.ParseFloat(os.Getenv("PERSONALIZATION_SCORE_WEIGHT"), 64)
	NORMALIZED_DISTANCE_WEIGHT, _ := strconv.ParseFloat(os.Getenv("NORMALIZED_DISTANCE_WEIGHT"), 64)

	for _, place := range originalScoresAndDistances {
		personalizationScore := 0.0
		for _, preference := range preferences {
			if preference == place.TypeOfDestination {
				personalizationScore = 1.0
				break
			}
		}
		normalizedDistance := (place.DistanceFromUser - minDistance) / (maxDistance - minDistance)
		newScore := VISIBILITY_SCORE_WEIGHT*place.VisitabilityScore + PERSONALIZATION_SCORE_WEIGHT*personalizationScore - NORMALIZED_DISTANCE_WEIGHT*normalizedDistance
		placeWithNewScore := place
		placeWithNewScore.VisitabilityScore = newScore
		placesWithIntermediateScores = append(placesWithIntermediateScores, placeWithNewScore)
	}
	return placesWithIntermediateScores
}

func (service *ItineraryService) getPlacesWithFinalScores(placesWithIntermediateScores []models.Place) (placesWithFinalScores []models.Place) {
	PENALTY_FACTOR_CONSTANT, _ := strconv.ParseFloat(os.Getenv("PENALTY_FACTOR_CONSTANT"), 64)
	penalties := make(map[string]int)
	for _, place := range placesWithIntermediateScores {
		penaltyFactor := penalties[place.TypeOfDestination]
		penalties[place.TypeOfDestination]++
		placeWithNewScore := place
		// **************** penalty constant (0.25) used here, will probably need to change ****************
		placeWithNewScore.VisitabilityScore = placeWithNewScore.VisitabilityScore - float64(penaltyFactor)*float64(penaltyFactor)*PENALTY_FACTOR_CONSTANT
		placesWithFinalScores = append(placesWithFinalScores, placeWithNewScore)
	}
	return placesWithFinalScores
}

func (service *ItineraryService) choosePlacesAccordingToNumberOfDays(placesWithFinalScores []models.Place, numberOfDaysAvailable int) (placesinItinerary []models.Place) {
	sort.Slice(
		placesWithFinalScores,
		func(i, j int) bool {
			return placesWithFinalScores[i].VisitabilityScore > placesWithFinalScores[j].VisitabilityScore
		},
	)
	if numberOfDaysAvailable*3 > len(placesWithFinalScores) {
		return placesWithFinalScores
	}
	return placesWithFinalScores[:3*numberOfDaysAvailable]
}

func (service *ItineraryService) setDayPlans(finalPlacesInItinerary []models.Place, numberOfDaysAvailable int, startLocationName, startLocationCity string, userStartLatitude float64, userStartLongitude float64) (finalDayPlans []models.DayPlan) {
	numberOfDaysInItinerary := min(numberOfDaysAvailable, len(finalPlacesInItinerary)/3+min(1, len(finalPlacesInItinerary)%3))
	for i := 0; i < numberOfDaysInItinerary; i++ {
		places := make([]models.Place, 0)
		places = append(places, models.Place{
			Name:      startLocationName,
			City:      startLocationCity,
			Latitude:  userStartLatitude,
			Longitude: userStartLongitude,
		})
		finalDayPlans = append(finalDayPlans, models.DayPlan{
			Day:    fmt.Sprintf("Day %v", i+1),
			Places: append(places, finalPlacesInItinerary[i*3:min(len(finalPlacesInItinerary), (i+1)*3)]...),
		})
	}
	return finalDayPlans
}
