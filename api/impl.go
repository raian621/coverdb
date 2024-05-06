package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/raian621/coverdb/database"
)

type Server struct{}

func (s Server) GetCoveragePath(w http.ResponseWriter, r *http.Request, path string, params GetCoveragePathParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) PostCoveragePath(w http.ResponseWriter, r *http.Request, path string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) GetHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func (s Server) DeleteCoveragePath(w http.ResponseWriter, r *http.Request, path string) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) DeleteKeys(w http.ResponseWriter, r *http.Request, params DeleteKeysParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) GetKeys(w http.ResponseWriter, r *http.Request, params GetKeysParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) PostKeys(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) PutKeys(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) PostUsersSignin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) PostUsersSignout(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (s Server) PostUsersSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userSigninData SigninDataBody

	if err := json.NewDecoder(r.Body).Decode(&userSigninData); err != nil {
		errorMsg := ErrorResponse{
			Message: "error in input data",
			Code:    http.StatusUnprocessableEntity,
		}
		if err := json.NewEncoder(w).Encode(errorMsg); err != nil {
			log.Printf("unexpected error: %v\n", err)
		}
	}

	if taken, err := database.UsernameTaken(userSigninData.Username); err != nil {
		errorMsg := ErrorResponse{
			Message: "error in input data",
			Code:    http.StatusUnprocessableEntity,
		}
		if err := json.NewEncoder(w).Encode(errorMsg); err != nil {
			log.Printf("unexpected error: %v\n", err)
		}
	} else if taken {
		errorMsg := ErrorResponse{
			Message: "username already exists",
			Code:    http.StatusConflict,
		}
		if err := json.NewEncoder(w).Encode(errorMsg); err != nil {
			log.Printf("unexpected error: %v\n", err)
		}
	}

	err := database.CreateUser(userSigninData.Username, userSigninData.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		errorMsg := ErrorResponse{
			Message: "username taken",
			Code:    http.StatusConflict,
		}
		if err := json.NewEncoder(w).Encode(errorMsg); err != nil {
			log.Printf("unexpected error: %v\n", err)
		}
	}

	response := SuccessfulSigninResponse{
		Message: "account successfully created",
		Code:    http.StatusCreated,
	}
	w.WriteHeader(response.Code)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("unexpected error: %v\n", err)
	}
}
