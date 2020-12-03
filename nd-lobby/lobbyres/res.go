package lobbyres

// StartMatchMakeRes is response of startmatchmake
type StartMatchMakeRes struct {
	TicketID string `json:"ticketid,string,omitempty"`
	ErrMsg   string `json:"errmsg,omitempty"`
}

// CancelMatchMakeRes is response of cancelmatchmake
type CancelMatchMakeRes struct {
	Status int32  `json:"status"`
	ErrMsg string `json:"errmsg,omitempty"`
}

// GetMatchMakeProcessRes is response of getmatchmakeprocess
type GetMatchMakeProcessRes struct {
	Status     int32  `json:"status"`
	Assignment string `json:"assignment,omitempty"`
	ErrMsg     string `json:"errmsg,omitempty"`
}
