package db

type MongoClient interface{
	Insert(user User) error
	FetchUserByID(id string) (*User, error)
	ListAllUsers() ([]*User, error)
	UpdateUser() error
	DeleteUser() error
	LoginUser() (*User, error)
	GetUserByEmail(email string) (*User, error)
}