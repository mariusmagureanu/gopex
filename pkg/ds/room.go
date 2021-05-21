package ds

import "gorm.io/gorm"

// Room represents the user for this application
//
// A user is the security principal for this application.
// It's also used as one of main axes for reporting.
//
// A user can have friends with whom they can share what they like.
//
// swagger:model
type Room struct {
	gorm.Model
	Name            string `json:"name"`
	Alias			string `json:"alias" gorm:"unique"`
	CostCenter      string `json:"cost_center"`
	GuestPin        string `json:"guest_pin"`
	HostPin         string `json:"host_pin"`
	Locked          bool   `json:"locked"`
	GuestMuted      bool   `json:"guest_muted"`
	AllowGuests     bool   `json:"allow_guests"`
	SSEon			bool   `json:"sse_on"`
}

func (Room) TableName() string {
	return "room"
}
