package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rigved-desai/paryatan-backend/api/interfaces"
	api "github.com/rigved-desai/paryatan-backend/api/utils"
)

type ItineraryController struct {
	ItineraryManager interfaces.ItineraryManager
}

func (controller *ItineraryController) GetItinerary(w http.ResponseWriter, r *http.Request) {

	var body struct {
		Location struct {
			Name      string  `json:"name"`
			City      string  `json:"city"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
		Preferences           []string `json:"preferences"`
		NumberOfDaysAvailable int      `json:"number_of_days_available"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		log.Println(err)
		api.RespondWithError(w, 500, "Error reading user input!")
		return
	}

	values, err := controller.ItineraryManager.GetItinerary(body.Location.Name, body.Location.City, body.Location.Latitude, body.Location.Longitude, body.Preferences, body.NumberOfDaysAvailable)
	if err != nil {
		log.Println(err)
		api.RespondWithError(w, 400, "Error getting itinerary!")
		return
	}

	api.RespondWithJSON(w, 200, values)
}
