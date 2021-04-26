package ds

type Layout struct {
	ID                int    `json:"id"`
	HostLayout        string `json:"host_layout"`
	GuestLayout       string `json:"guest_layout"`
	UseGuestLayout    bool   `json:"use_guest_layout"`
	PlusN             bool   `json:"plus_n"`
	ActorsOverlayText bool   `json:"actors_overlay_text"`
	FirstActor        string `json:"first_actor"`
}

func (Layout) TableName() string {
	return "layouts"
}
