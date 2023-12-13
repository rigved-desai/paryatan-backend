package models

type Itenary struct {
	DayPlans []DayPlan
}

type DayPlan struct {
	Day string
	Places []Place
}

type Place struct {
	Name string
	City string
	TypeOfDestination string
	DistanceFromUser float64
	Latitude float64
	Longitude float64
	Rating float64
	VisitabilityScore float64
}