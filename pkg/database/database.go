package database

import "fmt"

func NewApartmentDatabase() (*ApartmentDatabase, error) {
	db, err := ConnectDatabase()
	if err != nil {
		return nil, err
	}
	return &ApartmentDatabase{db}, err
}

func (db *ApartmentDatabase) Close() error {
	return db.Close()
}

func (db *ApartmentDatabase) CreateHouse(address, developer string, year int) (*House, error) {
	query := `INSERT INTO houses (address, developer, year)
VALUES ($1, $2, $3) RETURNING id, address, year, developer, created_at, updated_at;`
	house := &House{}
	err := db.QueryRow(query, address, developer, year).Scan(
		&house.Id,
		&house.Address,
		&house.Year,
		&house.Developer,
		&house.CreatedAt,
		&house.UpdateAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create house: %w", err)
	}
	return house, nil
}

func (db *ApartmentDatabase) UpdateHouse(house_id int) error {
	query := `UPDATE houses
SET updated_at = CURRENT_TIMESTAMP
WHERE id = $1;`
	_, err := db.Exec(query, house_id)
	if err != nil {
		return fmt.Errorf("failed to update house: %w", err)
	}
	return nil
}

func (db *ApartmentDatabase) CreateFlat(house_id, price int, rooms int) (*Flat, error) {
	query := `INSERT INTO flats (house_id, price, rooms, status)
VALUES ($1, $2, $3, 'created')
RETURNING id, house_id, price, rooms, status;`
	flat := &Flat{}
	err := db.QueryRow(query, house_id, price, rooms).Scan(
		&flat.Id,
		&flat.HouseId,
		&flat.Price,
		&flat.Rooms,
		&flat.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to create flat: %w", err)
	}
	db.UpdateHouse(house_id)
	return flat, nil
}

func (db *ApartmentDatabase) GetSubscribers(house_id int) ([]string, error) {
	query := `SELECT email
FROM subscribers
WHERE house_id = $1;`
	rows, err := db.Query(query, house_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscribers: %w", err)
	}

	defer rows.Close()
	subscribers := []string{}
	for rows.Next() {
		email := ""

		err = rows.Scan(&email)
		if err != nil {
			return nil, fmt.Errorf("failed to get subscribers: %w", err)
		}
		subscribers = append(subscribers, email)
	}
	return subscribers, nil
}

func (db *ApartmentDatabase) ModerateFlat(house_id, id int, status string) (*Flat, error) {
	query := `UPDATE flats
SET status = $3
WHERE house_id = $1 AND  id = $2
RETURNING id, house_id, price, rooms, status;`
	flat := &Flat{}
	err := db.QueryRow(query, house_id, id, status).Scan(
		&flat.Id,
		&flat.HouseId,
		&flat.Price,
		&flat.Rooms,
		&flat.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to moderate flat: %w", err)
	}
	db.UpdateHouse(house_id)
	return flat, nil
}

func (db *ApartmentDatabase) GetFlatList(house_id int, user_type string) ([]Flat, error) {
	query := ``
	if user_type == "client" {
		query = `SELECT f.house_id, f.id, f.price, f.rooms, f.status
FROM flats f
WHERE f.house_id = $1 AND f.status = 'approved'`
	} else {
		query = `SELECT f.house_id, f.id, f.price, f.rooms, f.status
FROM flats f
WHERE f.house_id = $1`
	}
	rows, err := db.Query(query, house_id)
	if err != nil {
		return nil, fmt.Errorf("failed to get flat list: %w", err)
	}

	defer rows.Close()
	flats := []Flat{}
	for rows.Next() {
		flat := &Flat{}

		err = rows.Scan(
			&flat.Id,
			&flat.HouseId,
			&flat.Price,
			&flat.Rooms,
			&flat.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to get flat list: %w", err)
		}
		flats = append(flats, *flat)
	}
	return flats, nil
}

func (db *ApartmentDatabase) AddSubscriber(house_id int, email string) error {
	query := `INSERT INTO subscribers (house_id, email)
VALUES ($1, $2);`
	_, err := db.Exec(query, house_id, email)
	if err != nil {
		return fmt.Errorf("failed to update house: %w", err)
	}
	return nil
}

func (db *ApartmentDatabase) AddUser(email, password, user_type string) (string, error) {
	query := `INSERT INTO users (email, password, type)
VALUES ($1, $2, $3) RETURNING UID;`
	uid := ""
	err := db.QueryRow(query, email, password, user_type).Scan(&uid)
	if err != nil {
		return "", fmt.Errorf("failed to add user: %w", err)
	}
	return uid, nil
}

func (db *ApartmentDatabase) GetUser(email string) (*User, error) {
	query := `SELECT *
FROM users
WHERE email = $1;`
	user := &User{}
	err := db.QueryRow(query, email).Scan(
		&user.UserId,
		&user.Email,
		&user.Password,
		&user.UserType)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}
