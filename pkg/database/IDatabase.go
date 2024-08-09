package database

import "database/sql"

type ApartmentDatabase struct {
	*sql.DB
}

type IApartmentStorage interface {
	Close() error

	CreateHouse(address, developer string, year int) (*House, error)
	UpdateHouse(house_id int) error

	CreateFlat(house_id, price int, rooms int) (*Flat, error)
	ModerateFlat(house_id, id int, status string) (*Flat, error)
	GetFlatList(house_id int, user_type string) ([]Flat, error)

	AddSubscriber(house_id int, email string) error
	GetSubscribers(house_id int) ([]string, error)

	AddUser(email, password, user_type string) (string, error)
	GetUser(email string) (*User, error)
}
