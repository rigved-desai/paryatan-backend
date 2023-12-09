package datastores

import (

	"github.com/rigved-desai/paryatan-backend/api/interfaces"
	"github.com/rigved-desai/paryatan-backend/api/models"
)

type ItenaryDataStore struct {
	interfaces.DBHandler
}

func (datastore *ItenaryDataStore) GetClusterByCoordinates(latitude, longitude float64) (cluster int, err error) {
	rows, err := datastore.Query("SELECT * FROM get_cluster($1, $2)", latitude, longitude)
	if err != nil {
		return 0, err
	}
	for rows.Next() {
		err = rows.Scan(&cluster)
		if err != nil {
			return 0, err
		}
	}
	return cluster, nil
}

func (datastore *ItenaryDataStore) GetPlacesWithOriginalScores(cluster int, latitude, longitude float64) (placesWithOriginalScores []models.Place, err error) {
	rows, err := datastore.Query("SELECT * FROM get_places_with_original_scores($1, $2, $3)", cluster, latitude, longitude)
	if err != nil {
		return []models.Place{}, err
	}
	for rows.Next() {
		var placeName, cityName, typeOfDestination string
		var latitude, longitude, rating, visitabilityScore, distanceFromUser float64
		err = rows.Scan(&placeName, &cityName, &typeOfDestination, &latitude, &longitude, &rating, &visitabilityScore, &distanceFromUser)
		if err != nil {
			return nil, err
		}

		placesWithOriginalScores = append(placesWithOriginalScores, models.Place{
			Name:              placeName,
			City:              cityName,
			TypeOfDestination: typeOfDestination,
			Latitude:          latitude,
			Longitude:         longitude,
			Rating:            rating,
			VisitabilityScore: visitabilityScore,
			DistanceFromUser: distanceFromUser,
		})
	}
	return placesWithOriginalScores, nil

}

func (datastore *ItenaryDataStore) GetMinAndMaxDistanceFromUser(cluster int, latitude, longitude float64) (minDistance float64, maxDistance float64, err error) {
	rows, err := datastore.Query("SELECT * FROM get_min_and_max_distance_from_user($1, $2, $3)", latitude, longitude, cluster)
	if err != nil {
		return 0.0, 0.0, err
	}

	for rows.Next() {

		err = rows.Scan(&minDistance, &maxDistance)
		if err != nil {
			return 0.0, 0.0, err
		}
	}
	return minDistance, maxDistance, nil
}

func (datastore *ItenaryDataStore) GetAllPlaces() ([]models.Place, error) {
	rows, err := datastore.Query("SELECT place_name, latitude_coordinates, longitude_coordinates, rating FROM tourist_places")
	if err != nil {
		return nil, err
	}
	var places []models.Place
	for rows.Next() {
		var place_name string
		var latitude, longitude, rating float64
		err = rows.Scan(&place_name, &latitude, &longitude, &rating)
		if err != nil {
			return nil, err
		}
		places = append(places, models.Place{
			Name:      place_name,
			Latitude:  latitude,
			Longitude: longitude,
			Rating:    rating,
		})
	}
	return places, nil

}