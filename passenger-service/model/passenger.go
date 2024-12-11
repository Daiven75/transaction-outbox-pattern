package model

import (
	"gorm.io/gorm"
)

type Passenger struct {
	gorm.Model
	FirstName string `json:"first_name"`
	PlanType  string `json:"plan_type"`
	Dispatch  bool   `json:"dispatch"`
	FlightID  uint   `json:"flight_id"`
}

type PassengerOutbox struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	FirstName string `json:"first_name"`
	FlightID  uint   `json:"flight_id"`
}

func (*PassengerOutbox) TableName() string {
	return "passenger_outbox"
}

func (p *PassengerOutbox) FromPassenger(passenger Passenger) {
	p.ID = passenger.ID
	p.FirstName = passenger.FirstName
	p.FlightID = passenger.FlightID
}
