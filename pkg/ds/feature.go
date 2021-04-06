// Package ds acts as a placeholder for all the
// model objects that are to be persisted.
package ds

type Feature struct {
	ID               int  `json:"id"`
	BasicView        bool `json:"basic_view"`
	SimpleView       bool `json:"simple_view"`
	DevelView        bool `json:"devel_view"`
	WebrtcView       bool `json:"webrtc_view"`
	Pin              bool `json:"pin"`
	Streaming        bool `json:"streaming"`
	Recording        bool `json:"recording"`
	DisposableAlias  bool `json:"disposable_alias"`
	Registrar        bool `json:"registrar"`
	Lecture          bool `json:"lecture"`
	C2bRoom          bool `json:"c2b_room"`
	AutomaticDialout bool `json:"automatic_dialout"`
	Allhands         bool `json:"allhands"`
	BoardMeeting     bool `json:"boardmeeting"`
	CameraCrew       bool `json:"camera_crew"`
	EventRoom        bool `json:"event_room"`
	Recents          bool `json:"recents"`
	CourtRoom        bool `json:"courtroom"`
	PrivateRoom      bool `json:"private_room"`
	ScreeningRoom    bool `json:"screening_room"`
	WaitingRoom      bool `json:"waiting_room"`
}

func (Feature) TableName() string {
	return "rooms_featuremodel"
}
