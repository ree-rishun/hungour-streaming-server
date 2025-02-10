package models

import (
	"time"
)

type Message struct {
	Role			string		`firestore:"role"`
	Text			string		`firestore:"text"`
}

type Process struct {
	Messages		[]Message	`firestore:"messages"`
	ConciergeId		string		`firestore:"concierge_id"`
	Status			string		`firestore:"status"`
	ReservedTime	time.Time	`firestore:"reserved_time"`
	CreatedAt		time.Time	`firestore:"created_at"`
	UpdatedAt		time.Time	`firestore:"updated_at"`
}
