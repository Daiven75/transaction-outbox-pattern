package model

import "gorm.io/gorm"

type Flight struct {
	gorm.Model
	Company     string      `json:"company"`
	Origin      string      `json:"origin"`
	Destination string      `json:"destination"`
	Passengers  []Passenger `json:"passengers" gorm:"foreignKey:FlightID"`
}

type Passenger struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	FlightID  uint   `json:"flight_id"`
}
