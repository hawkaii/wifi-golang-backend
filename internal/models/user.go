package models

type User struct {
	ID    string `bson:"_id,omitempty"`
	Email string `bson:"email"`
	Name  string `bson:"name"`
}
