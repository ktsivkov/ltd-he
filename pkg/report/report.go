package report

type Report struct {
	File       string `json:"file"`
	LastGameId int    `json:"lastGameId"`
	Token      string `json:"token"`
	Payload    []byte `json:"payload"`
}
