package node

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	PENDING Status = "pending"
	ONLINE  Status = "online"
	OFFLINE Status = "offline"
)

type Node struct {
	ID              uuid.UUID
	OrganizationID  uuid.UUID
	Name            string
	Hostname        string
	Status          Status
	OperatingSystem string
	AgentVersion    string
	LastActive      time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type PairingToken struct {
	NodeID    uuid.UUID
	Token     string
	ExpiresAt time.Time
}
