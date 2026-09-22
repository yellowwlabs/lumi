package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"lumi.yellowlabs.space/internal/auth"
	"lumi.yellowlabs.space/internal/node"
	"lumi.yellowlabs.space/internal/organizations"
)

type createNodeRequest struct {
	Name string `json:"name"`
}

type createNodeResponse struct {
	Node         *node.Node `json:"node"`
	PairingToken string     `json:"pairing_token"`
}

type updateNodeRequest struct {
	Name            string `json:"name"`
	Hostname        string `json:"hostname"`
	OperatingSystem string `json:"operating_system"`
	AgentVersion    string `json:"agent_version"`
}

type registerAgentRequest struct {
	PairingToken    string `json:"pairing_token"`
	Hostname        string `json:"hostname"`
	OperatingSystem string `json:"operating_system"`
	AgentVersion    string `json:"agent_version"`
}

type registerAgentResponse struct {
	NodeID     string `json:"node_id"`
	AgentToken string `json:"agent_token"`
}

func nodeErrorStatus(err error) (int, string) {
	switch err {
	case organizations.ErrForbidden:
		return http.StatusForbidden, "not a member of this organization"
	case node.ErrNodeNotFound:
		return http.StatusNotFound, "node not found"
	default:
		return http.StatusInternalServerError, "internal server error"
	}
}

func SetupNodeRoutes(router chi.Router, nodeService *node.NodeService) {
	router.Route("/organizations/{organizationID}/nodes", func(r chi.Router) {
		r.Use(auth.RequireAuth)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			organizationID := chi.URLParam(r, "organizationID")
			requesterID, _ := auth.UserIDFromContext(r.Context())

			nodes, err := nodeService.GetNodesByOrganizationID(organizationID, requesterID)
			if err != nil {
				status, msg := nodeErrorStatus(err)
				writeError(w, status, msg)
				return
			}

			writeJSON(w, http.StatusOK, nodes)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			organizationID := chi.URLParam(r, "organizationID")
			requesterID, _ := auth.UserIDFromContext(r.Context())

			var req createNodeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
				writeError(w, http.StatusBadRequest, "name is required")
				return
			}

			n, pairingToken, err := nodeService.CreateNode(organizationID, requesterID, req.Name)
			if err != nil {
				status, msg := nodeErrorStatus(err)
				writeError(w, status, msg)
				return
			}

			writeJSON(w, http.StatusCreated, createNodeResponse{Node: n, PairingToken: pairingToken})
		})

		r.Get("/{nodeID}", func(w http.ResponseWriter, r *http.Request) {
			organizationID := chi.URLParam(r, "organizationID")
			nodeID := chi.URLParam(r, "nodeID")
			requesterID, _ := auth.UserIDFromContext(r.Context())

			n, err := nodeService.GetNode(organizationID, requesterID, nodeID)
			if err != nil {
				status, msg := nodeErrorStatus(err)
				writeError(w, status, msg)
				return
			}

			writeJSON(w, http.StatusOK, n)
		})

		r.Put("/{nodeID}", func(w http.ResponseWriter, r *http.Request) {
			organizationID := chi.URLParam(r, "organizationID")
			nodeID := chi.URLParam(r, "nodeID")
			requesterID, _ := auth.UserIDFromContext(r.Context())

			var req updateNodeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, http.StatusBadRequest, "invalid request payload")
				return
			}

			n, err := nodeService.UpdateNode(organizationID, requesterID, nodeID, req.Name, req.Hostname, req.OperatingSystem, req.AgentVersion)
			if err != nil {
				status, msg := nodeErrorStatus(err)
				writeError(w, status, msg)
				return
			}

			writeJSON(w, http.StatusOK, n)
		})

		r.Delete("/{nodeID}", func(w http.ResponseWriter, r *http.Request) {
			organizationID := chi.URLParam(r, "organizationID")
			nodeID := chi.URLParam(r, "nodeID")
			requesterID, _ := auth.UserIDFromContext(r.Context())

			if err := nodeService.DeleteNode(organizationID, requesterID, nodeID); err != nil {
				status, msg := nodeErrorStatus(err)
				writeError(w, status, msg)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})
	})
}

// SetupAgentRoutes wires the endpoints called by the agent process running on
// a paired server, not by an authenticated dashboard user. Auth here is the
// pairing/agent token itself, not a user session.
func SetupAgentRoutes(router chi.Router, nodeService *node.NodeService) {
	router.Route("/agent", func(r chi.Router) {
		r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
			var req registerAgentRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PairingToken == "" {
				writeError(w, http.StatusBadRequest, "pairing_token is required")
				return
			}

			n, agentToken, err := nodeService.RegisterAgent(req.PairingToken, req.Hostname, req.OperatingSystem, req.AgentVersion)
			if err != nil {
				switch err {
				case node.ErrPairingTokenInvalid, node.ErrPairingTokenExpired:
					writeError(w, http.StatusUnauthorized, "pairing token invalid or expired")
				default:
					writeError(w, http.StatusInternalServerError, "internal server error")
				}
				return
			}

			writeJSON(w, http.StatusOK, registerAgentResponse{NodeID: n.ID.String(), AgentToken: agentToken})
		})

		r.Post("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
			agentToken := r.Header.Get("X-Agent-Token")
			if agentToken == "" {
				writeError(w, http.StatusUnauthorized, "X-Agent-Token header is required")
				return
			}

			if err := nodeService.Heartbeat(agentToken); err != nil {
				writeError(w, http.StatusUnauthorized, "invalid agent token")
				return
			}

			w.WriteHeader(http.StatusNoContent)
		})
	})
}
