package report

type Report struct {
	LastGameId int    `json:"lastGameId"`
	Token      string `json:"token"`
	Payload    []byte `json:"payload"`
}
