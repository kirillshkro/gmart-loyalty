package services

import "net/http"

type Registerer interface {
	Register(w http.ResponseWriter, r *http.Request)
}

func (u UserService) Register(w http.ResponseWriter, r *http.Request) {
	// Implement user registration logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User registered successfully"))
}
