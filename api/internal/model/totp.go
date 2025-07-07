package model

type TOTPNew struct {
	Secret  string `json:"secret"`
	Account string `json:"account"`
	URI     string `json:"uri"`
}

type TOTPBackup struct {
	Backup string `json:"backup"`
}
