package models

type User struct {
	ReserveName		string		`firestore:"reserve_name"`
	LineId			string		`firestore:"line_id"`
	Tel				string		`firestore:"tel"`
}
