package ds

type Room struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	CostCenter      string `json:"cost_center"`
	Locked          bool   `json:"locked"`
	GuestMuted      bool   `json:"guest_muted"`
	UUID            string `json:"uuid"`
	VoipPin         string `json:"voip_pin"`
	AliasID         int    `json:"alias_id"`
	IvrThemeID      int    `json:"ivr_theme_id"`
	LayoutCommandID int    `json:"layout_command_id"`
	TenantID        int    `json:"tenant_id"`
	FeatureID       int    `json:"feature_id"`
	GuestPin        string `json:"guest_pin"`
	HostPin         string `json:"host_pin"`
	AllowGuests     bool   `json:"allow_guests"`
}

func (Room) TableName() string {
	return "rooms"
}
