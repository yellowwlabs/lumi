package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"lumi.yellowlabs.space/internal/auth"
)

type registerRequest struct {
	DisplayName string `json:"display_name"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

type resetPasswordRequest struct {
	ResetToken  string `json:"reset_token"`
	NewPassword string `json:"new_password"`
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func SetupAuthRoutes(router chi.Router, userService *auth.UserService) {
	router.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request payload")
			return
		}
		if req.Username == "" || req.Email == "" || req.Password == "" {
			writeError(w, http.StatusBadRequest, "username, email and password are required")
			return
		}

		user := &auth.User{
			DisplayName: req.DisplayName,
			Username:    req.Username,
			Email:       req.Email,
		}

		if err := userService.RegisterUser(user, req.Password); err != nil {
			if err == auth.ErrUserAlreadyExists {
				writeError(w, http.StatusConflict, "user already exists")
				return
			}
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		writeJSON(w, http.StatusCreated, map[string]string{
			"id":       user.ID.String(),
			"username": user.Username,
		})
	})

	router.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request payload")
			return
		}

		user, accessToken, refreshToken, err := userService.Login(req.Username, req.Password, r.RemoteAddr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"user_id":       user.ID.String(),
		})
	})

	router.Post("/refresh", func(w http.ResponseWriter, r *http.Request) {
		var req refreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
			writeError(w, http.StatusBadRequest, "refresh_token is required")
			return
		}

		accessToken, refreshToken, err := userService.Refresh(req.RefreshToken, r.RemoteAddr)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired refresh token")
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	})

	router.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		var req logoutRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
			writeError(w, http.StatusBadRequest, "refresh_token is required")
			return
		}

		if err := userService.Logout(req.RefreshToken); err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	router.Post("/forgot-password", func(w http.ResponseWriter, r *http.Request) {
		var req forgotPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
			writeError(w, http.StatusBadRequest, "email is required")
			return
		}

		resetToken, err := userService.ForgotPassword(req.Email, r.RemoteAddr)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// TODO: email resetToken to the user instead of returning it once
		// mail delivery is wired up. Returned directly here since there's no
		// email provider configured yet.
		resp := map[string]string{"message": "if the account exists, a reset token has been issued"}
		if resetToken != "" {
			resp["reset_token"] = resetToken
		}
		writeJSON(w, http.StatusOK, resp)
	})

	router.Post("/reset-password", func(w http.ResponseWriter, r *http.Request) {
		var req resetPasswordRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ResetToken == "" || req.NewPassword == "" {
			writeError(w, http.StatusBadRequest, "reset_token and new_password are required")
			return
		}

		if err := userService.ResetPassword(req.ResetToken, req.NewPassword); err != nil {
			writeError(w, http.StatusBadRequest, "invalid or expired reset token")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	router.With(auth.RequireAuth).Get("/@me", func(w http.ResponseWriter, r *http.Request) {
		userID, _ := auth.UserIDFromContext(r.Context())
		writeJSON(w, http.StatusOK, map[string]string{"user_id": userID})
	})
}
