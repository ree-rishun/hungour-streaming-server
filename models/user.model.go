package models

type User struct {
	ReserveName		string		`firestore:"reserve_name"`
	Tel				string		`firestore:"tel"`
}
