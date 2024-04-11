package api

type SuccessRes struct {
	Message string `json:"message"`
}

type ErrorRes struct {
	Error string `json:"error"`
}
