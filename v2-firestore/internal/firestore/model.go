package firestore

import "time"

type Todo struct {
	ID        string    `firestore:"id,omitempty"`
	UserID    string    `firestore:"userId,omitempty"`
	Title     string    `firestore:"title,omitempty"`
	Done      bool      `firestore:"done"`
	CreatedAt time.Time `firestore:"createdAt,omitempty"`
}
