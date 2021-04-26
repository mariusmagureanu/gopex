package ds

type User struct {
	ID              int
	Email           string
	FirstName       string
	LastName        string
	IsSuperUser     bool
	IsRoomUser      bool
	IsRoomSuperUser bool
	IsStaff         bool
	IsActive        bool
	IsAdmin         bool
}

func (User) TableName() string {
	return "users"
}
