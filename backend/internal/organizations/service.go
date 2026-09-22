package organizations

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationService struct {
	repo *OrganizationRepository
}

func NewOrganizationService(repo *OrganizationRepository) *OrganizationService {
	return &OrganizationService{repo: repo}
}

func isValidRole(role MemberRole) bool {
	switch role {
	case ADMIN, MEMBER, OWNER:
		return true
	default:
		return false
	}
}

func canManageMembers(role MemberRole) bool {
	return role == ADMIN || role == OWNER
}

// GetMember returns the caller's own membership record, used for authorization checks.
func (s *OrganizationService) GetMember(organizationID, userID string) (*OrganizationMembers, error) {
	return s.repo.GetOrganizationMember(organizationID, userID)
}

func (s *OrganizationService) GetOrganizationMembers(organizationID, requesterID string) ([]*OrganizationMembers, error) {
	requester, err := s.repo.GetOrganizationMember(organizationID, requesterID)
	if err != nil {
		return nil, err
	}
	if requester == nil {
		return nil, ErrForbidden
	}

	return s.repo.GetOrganizationMembers(organizationID)
}

func (s *OrganizationService) InviteUserToOrganization(organizationID, requesterID, userID, role string) error {
	requester, err := s.repo.GetOrganizationMember(organizationID, requesterID)
	if err != nil {
		return err
	}
	if requester == nil || !canManageMembers(requester.Role) {
		return ErrForbidden
	}

	memberRole := MemberRole(role)
	if !isValidRole(memberRole) {
		return ErrInvalidRole
	}

	parsedOrganizationID, err := uuid.Parse(organizationID)
	if err != nil {
		return err
	}
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	member := &OrganizationMembers{
		OrganizationID: parsedOrganizationID,
		UserID:         parsedUserID,
		Role:           memberRole,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	member.InvitationStatus = PENDING

	return s.repo.CreateOrganizationMember(member)
}
