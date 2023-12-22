package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/rigved-desai/paryatan-backend/api/controllers"
	"github.com/rigved-desai/paryatan-backend/api/datastores"
	"github.com/rigved-desai/paryatan-backend/api/services"
	"github.com/rigved-desai/paryatan-backend/db"
)

func NewRouter(postgre *db.Postgre) *chi.Mux {
	router := chi.NewRouter()
	router.Mount("/v1", v1Router(postgre))
	return router
}

func v1Router(postgre *db.Postgre) *chi.Mux {
	router := chi.NewRouter()

	postgreSQLHandler := &datastores.PostgreSQLHanlder{
		ConnPool: postgre.DB,
	}

	controller := controllers.ItineraryController{
		ItineraryManager: &services.ItineraryService{
			DataAccessor: &datastores.ItineraryDataStore{
				DBHandler: postgreSQLHandler,
			},
		},
	}
	router.Post("/itinerary", controller.GetItinerary)
	return router
}
