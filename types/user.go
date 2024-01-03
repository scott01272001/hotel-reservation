package types

type CreateUserparam struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type User struct {
	ID                string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string `bson:"firstName" json:"firstName"`
	LastName          string `bson:"lasttName" json:"lastName"`
	Email             string `bson:"email" json:"email"`
	EncryptedPasswrod string `bson:"envryptedPassword" json:"-"`
}
