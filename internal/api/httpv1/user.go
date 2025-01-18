package httpv1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tribe-payments-wallet-golang-interview-assignment/internal/auth"
	"tribe-payments-wallet-golang-interview-assignment/internal/config"
	"tribe-payments-wallet-golang-interview-assignment/internal/user"

	"github.com/sumup-oss/go-pkgs/logger"
)

func HandlerCreateUser(log logger.StructuredLogger, userService *user.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		log.Info("HandlerCreateUser")

		var req user.UserCreateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error(config.ErrInvalidRequestPayload.Error())
			http.Error(w, config.ErrInvalidRequestPayload.Error(), http.StatusBadRequest)
			return
		}

		log.Info(fmt.Sprintf("Creating user with email: %s", req.Email))

		usr := &user.User{
			Email:        req.Email,
			PasswordHash: req.Password,
			FirstName:    req.FirstName,
			LastName:     req.LastName,
		}

		err = userService.CreateUser(r.Context(), usr)
		if err != nil {
			log.Error(fmt.Sprintf(config.ErrFailedToCreateUser.Error(), err))
			http.Error(w, config.ErrFailedToCreateUser.Error(), http.StatusUnprocessableEntity)
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
			log.Error(config.ErrInvalidRequestPayload.Error())
			http.Error(w, config.ErrInvalidRequestPayload.Error(), http.StatusBadRequest)
			return
		}

		usr, err := userService.GetUserByEmail(r.Context(), req.Email)

		err = auth.CheckPasswordHash(req.Password, usr.PasswordHash)
		if err != nil {
			log.Error(config.ErrInvalidCredentials.Error())
			http.Error(w, config.ErrInvalidCredentials.Error(), http.StatusUnauthorized)
			return
		}

		jwt, err := auth.GenerateJWT(usr.ID.String(), usr.Email)
		if err != nil {
			log.Error(fmt.Sprintf(config.ErrGeneratingJWT.Error(), err))
			http.Error(w, config.ErrGeneratingJWT.Error(), http.StatusUnprocessableEntity)
			return
		}

		// Return the JWT token in the response
		response := map[string]string{
			"token": jwt,
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error(fmt.Sprintf("Error encoding response: %s", err))
			http.Error(w, fmt.Sprintf("Error encoding response: %s", err), http.StatusUnprocessableEntity)
			return
		}
	}
}
