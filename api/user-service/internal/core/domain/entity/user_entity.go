package entity

type UserEntity struct {
	ID         int64
	Name       string
	Email      string
	Password   string
	RoleName   string
	RoleId     int64
	Address    string
	Lat        string
	Lng        string
	Phone      string
	Photo      string
	IsVerified bool
	Token      string
}

type QueryStringEntity struct {
	Search    string
	Page      int64
	Limit     int64
	OrderBy   string
	OrderType string
}
