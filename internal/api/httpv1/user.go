package httpv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tribe-payments-wallet-golang-interview-assignment/internal/auth"
	"tribe-payments-wallet-golang-interview-assignment/internal/user"

	"github.com/sumup-oss/go-pkgs/errors"
	"github.com/sumup-oss/go-pkgs/logger"
)

var (
	ErrInvalidRequestPayload = errors.New("Invalid request payload")
	ErrInvalidCredentials    = errors.New("Invalid credentials")
)

func HandlerCreateUser(log logger.StructuredLogger, userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Info("HandlerCreateUser")

		var req user.UserCreateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error(ErrInvalidRequestPayload.Error())
			http.Error(w, ErrInvalidRequestPayload.Error(), http.StatusBadRequest)
			return
		}

		usr := &user.User{
			Email:        req.Email,
			PasswordHash: req.Password,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}

		err = userService.CreateUser(r.Context(), usr)
		if err != nil {
			log.Error(fmt.Sprintf("Error creating user: %s", err))
			http.Error(w, fmt.Sprintf("Error creating user: %s", err), http.StatusUnprocessableEntity)
			return
		}
	}
}

func HandlerLoginUser(log logger.StructuredLogger, userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		log.Info("Login")

		var req user.UserLoginRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error(ErrInvalidRequestPayload.Error())
			http.Error(w, ErrInvalidRequestPayload.Error(), http.StatusBadRequest)
			return
		}

		usr, err := userService.GetUserByEmail(r.Context(), req.Email)

		err = auth.CheckPasswordHash(req.Password, usr.PasswordHash)
		if err != nil {
			log.Error(ErrInvalidCredentials.Error())
			http.Error(w, ErrInvalidCredentials.Error(), http.StatusUnauthorized)
			return
		}

		jwt, err := auth.GenerateJWT(usr.ID.String(), usr.Email)
		if err != nil {
			log.Error(fmt.Sprintf("Error generating JWT: %s", err))
			http.Error(w, fmt.Sprintf("Error generating JWT: %s", err), http.StatusInternalServerError)
			return
		}

		// Return the JWT token in the response
		response := map[string]string{
			"token": jwt,
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error(fmt.Sprintf("Error encoding response: %s", err))
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusInternalServerError)
			return
		}
	}
}
