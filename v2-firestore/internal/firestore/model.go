package firestore

type Todo struct {
	ID     string `firestore:"id,omitempty"`
	UserID string `firestore:"userId,omitempty"`
	Title  string `firestore:"title,omitempty"`
}
