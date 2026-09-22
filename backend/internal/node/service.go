package node

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"lumi.yellowlabs.space/internal/organizations"
)

const pairingTokenTTL = 15 * time.Minute

type NodeService struct {
	repo *NodeRepository
	orgs *organizations.OrganizationService
}

func NewNodeService(repo *NodeRepository, orgs *organizations.OrganizationService) *NodeService {
	return &NodeService{repo: repo, orgs: orgs}
}

func (s *NodeService) requireMember(organizationID, requesterID string) error {
	member, err := s.orgs.GetMember(organizationID, requesterID)
	if err != nil {
		return err
	}
	if member == nil {
		return organizations.ErrForbidden
	}
	return nil
}

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *NodeService) GetNodesByOrganizationID(organizationID, requesterID string) ([]*Node, error) {
	if err := s.requireMember(organizationID, requesterID); err != nil {
		return nil, err
	}
	return s.repo.GetNodesByOrganizationID(organizationID)
}

func (s *NodeService) GetNode(organizationID, requesterID, nodeID string) (*Node, error) {
	if err := s.requireMember(organizationID, requesterID); err != nil {
		return nil, err
	}

	node, err := s.repo.GetNodeByID(nodeID)
	if err != nil {
		return nil, err
	}
	if node == nil || node.OrganizationID.String() != organizationID {
		return nil, ErrNodeNotFound
	}

	return node, nil
}

// CreateNode registers a placeholder node for the org and returns a one-time
// pairing token. The plaintext token is never stored; only its hash is.
func (s *NodeService) CreateNode(organizationID, requesterID, name string) (*Node, string, error) {
	if err := s.requireMember(organizationID, requesterID); err != nil {
		return nil, "", err
	}

	parsedOrganizationID, err := uuid.Parse(organizationID)
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	n := &Node{
		ID:             uuid.New(),
		OrganizationID: parsedOrganizationID,
		Name:           name,
		Status:         PENDING,
		LastActive:     now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.CreateNode(n); err != nil {
		return nil, "", err
	}

	pairingToken, err := randomToken()
	if err != nil {
		return nil, "", err
	}
	if err := s.repo.SetPairingToken(n.ID.String(), hashToken(pairingToken), now.Add(pairingTokenTTL)); err != nil {
		return nil, "", err
	}

	return n, pairingToken, nil
}

func (s *NodeService) UpdateNode(organizationID, requesterID, nodeID, name, hostname, operatingSystem, agentVersion string) (*Node, error) {
	node, err := s.GetNode(organizationID, requesterID, nodeID)
	if err != nil {
		return nil, err
	}

	node.Name = name
	node.Hostname = hostname
	node.OperatingSystem = operatingSystem
	node.AgentVersion = agentVersion
	node.UpdatedAt = time.Now()

	if err := s.repo.UpdateNode(node); err != nil {
		return nil, err
	}

	return node, nil
}

func (s *NodeService) DeleteNode(organizationID, requesterID, nodeID string) error {
	if _, err := s.GetNode(organizationID, requesterID, nodeID); err != nil {
		return err
	}

	return s.repo.DeleteNode(nodeID)
}

// RegisterAgent is called by the agent running on the target server, not by
// an authenticated dashboard user. It exchanges a short-lived pairing token
// for a long-lived agent token used for subsequent heartbeats.
func (s *NodeService) RegisterAgent(pairingToken, hostname, operatingSystem, agentVersion string) (*Node, string, error) {
	node, expiresAt, err := s.repo.GetNodeByPairingTokenHash(hashToken(pairingToken))
	if err != nil {
		return nil, "", err
	}
	if node == nil {
		return nil, "", ErrPairingTokenInvalid
	}
	if expiresAt == nil || time.Now().After(*expiresAt) {
		return nil, "", ErrPairingTokenExpired
	}

	agentToken, err := randomToken()
	if err != nil {
		return nil, "", err
	}

	now := time.Now()
	if err := s.repo.CompletePairing(node.ID.String(), hashToken(agentToken), hostname, operatingSystem, agentVersion, now); err != nil {
		return nil, "", err
	}

	node.Hostname = hostname
	node.OperatingSystem = operatingSystem
	node.AgentVersion = agentVersion
	node.Status = ONLINE
	node.LastActive = now

	return node, agentToken, nil
}

// Heartbeat is called periodically by a connected agent to prove liveness.
func (s *NodeService) Heartbeat(agentToken string) error {
	node, err := s.repo.GetNodeByAgentTokenHash(hashToken(agentToken))
	if err != nil {
		return err
	}
	if node == nil {
		return ErrAgentTokenInvalid
	}

	return s.repo.UpdateNodeStatus(node.ID.String(), ONLINE, time.Now())
}
