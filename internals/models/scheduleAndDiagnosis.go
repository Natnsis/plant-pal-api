package models

import "gorm.io/gorm"

type Schedule struct {
	gorm.Model
	Name       string
	Duration   string
	CarePlanID uint
}

type Diagnosis struct {
	gorm.Model
	ImageUrl        string
	Symptoms        []string
	DiagnosisResult string
}
