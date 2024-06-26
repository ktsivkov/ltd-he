package player

type Player struct {
	BattleTag              string `json:"battleTag"`
	LogsPathAbsolute       string `json:"logsPathAbsolute"`
	LogsPathRelative       string `json:"logsPathRelative"`
	ReportFilePathAbsolute string `json:"reportFilePathAbsolute"`
	ReportFilePathRelative string `json:"reportFilePathRelative"`
}
