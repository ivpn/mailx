package model

type TOTPNew struct {
	Secret string `json:"secret"`
	URI    string `json:"uri"`
}

type TOTPBackup struct {
	Backup string `json:"backup"`
}
