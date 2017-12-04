package model

//AuthServiceError structure contains the error resource definition
type AuthServiceError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
