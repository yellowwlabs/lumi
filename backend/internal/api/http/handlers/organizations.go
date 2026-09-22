package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"lumi.yellowlabs.space/internal/auth"
	"lumi.yellowlabs.space/internal/organizations"
)

type inviteMemberRequest struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func SetupOrganizationRoutes(router chi.Router, orgService *organizations.OrganizationService) {
	router.Route("/organizations/{organizationID}", func(r chi.Router) {
		r.Use(auth.RequireAuth)

		r.Get("/members", func(w http.ResponseWriter, r *http.Request) {
			organizationID := chi.URLParam(r, "organizationID")
			requesterID, _ := auth.UserIDFromContext(r.Context())

			members, err := orgService.GetOrganizationMembers(organizationID, requesterID)
			if err != nil {
				if err == organizations.ErrForbidden {
					writeError(w, http.StatusForbidden, "not a member of this organization")
					return
				}
				writeError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			writeJSON(w, http.StatusOK, members)
		})

		r.Post("/members", func(w http.ResponseWriter, r *http.Request) {
			organizationID := chi.URLParam(r, "organizationID")
			requesterID, _ := auth.UserIDFromContext(r.Context())

			var req inviteMemberRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.Role == "" {
				writeError(w, http.StatusBadRequest, "user_id and role are required")
				return
			}

			if err := orgService.InviteUserToOrganization(organizationID, requesterID, req.UserID, req.Role); err != nil {
				switch err {
				case organizations.ErrForbidden:
					writeError(w, http.StatusForbidden, "not authorized to invite members")
				case organizations.ErrInvalidRole:
					writeError(w, http.StatusBadRequest, "invalid role")
				default:
					writeError(w, http.StatusInternalServerError, "internal server error")
				}
				return
			}

			w.WriteHeader(http.StatusCreated)
		})
	})
}
