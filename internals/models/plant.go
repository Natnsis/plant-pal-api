package models

import "gorm.io/gorm"

type Watering struct {
	gorm.Model
	Frequency string
	Amount    string
	Method    string
	Tips      string
}

type Light struct {
	Requirment string
	Duration   string
	Placement  string
	Rotate     string
}

type Humidity struct {
	IdealRange string
	BoostTips  string
	Frequency  string
}

type Tempreture struct {
	IdealRange   string
	MinTolerance string
	Avoid        string
}

type Soil struct {
	SoilType        string
	PhRange         string
	ReportFrequency string
	NextReport      string
}

type Fertilizar struct {
	FertilizerType string
	Season         string
	Dilution       string
}

type Pruning struct {
	When   string
	Remove string
	Tools  string
}

type CarePlan struct {
	gorm.Model
	Watering   Watering
	Light      Light
	Humidity   Humidity
	Tempreture Tempreture
	Soil       Soil
	Fertilizar Fertilizar
	Pruning    Pruning
}

type Plant struct {
	gorm.Model
	Name           string
	ScientificName string
	Family         string
	Origin         string
	Location       string
	IsSafe         string
	Difficulty     string
	ImageUrl       string
}
