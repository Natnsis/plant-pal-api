package models

type Analytic struct {
	gorm.Model
	PostCount int
	Streak    int
	Tasks     int
	Entries   int
}
