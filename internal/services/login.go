package services

import "net/http"

type Loginer interface {
	Login(w http.ResponseWriter, r *http.Request)
}

func (u UserService) Login(w http.ResponseWriter, r *http.Request) {
	// Implement user login logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User logged in successfully"))
}
