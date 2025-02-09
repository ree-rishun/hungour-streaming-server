package models

import(
	"time"
)

type Concierge struct {
	ReserveList		[]ReserveList	`firestore:"reserve_list"`
	UserId			string			`firestore:"user_id"`
	Status			string			`firestore:"status"`
	SeatType		string			`firestore:"seat_type"`
	PartySize		int64			`firestore:"party_size"`
	DepartureTime	int64			`firestore:"departure_time"`
	Cursor			int64			`firestore:"cursor"`
	CreatedAt		time.Time		`firestore:"created_at"`
	UpdatedAt		time.Time		`firestore:"updated_at"`
}

type ReserveList struct {
	Id				string			`firestore:"id"`
	Name			string			`firestore:"name"`
	Tel				string			`firestore:"tel"`
}
