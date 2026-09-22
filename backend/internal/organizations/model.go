package organizations

import (
	"time"

	"github.com/google/uuid"
)

type MemberRole string

type InvitationStatus string

const (
	ADMIN  MemberRole = "admin"
	MEMBER MemberRole = "member"
	OWNER  MemberRole = "owner"
)

const (
	ACCEPTED InvitationStatus = "accepted"
	DECLINED InvitationStatus = "declined"
	PENDING  InvitationStatus = "pending"
)

type Organization struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OrganizationMembers struct {
	OrganizationID   uuid.UUID
	UserID           uuid.UUID
	Role             MemberRole
	InvitationStatus InvitationStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
