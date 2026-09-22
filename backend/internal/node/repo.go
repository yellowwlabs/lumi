package node

import (
	"database/sql"
	"time"
)

type NodeRepository struct {
	db *sql.DB
}

func NewNodeRepository(db *sql.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) GetNodesByOrganizationID(organizationID string) ([]*Node, error) {
	rows, err := r.db.Query("SELECT id, organization_id, name, hostname, status, operating_system, agent_version, last_active, created_at, updated_at FROM nodes WHERE organization_id = $1", organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*Node
	for rows.Next() {
		var node Node
		if err := rows.Scan(&node.ID, &node.OrganizationID, &node.Name, &node.Hostname, &node.Status, &node.OperatingSystem, &node.AgentVersion, &node.LastActive, &node.CreatedAt, &node.UpdatedAt); err != nil {
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *NodeRepository) CreateNode(node *Node) error {
	_, err := r.db.Exec("INSERT INTO nodes (id, organization_id, name, hostname, status, operating_system, agent_version, last_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		node.ID, node.OrganizationID, node.Name, node.Hostname, node.Status, node.OperatingSystem, node.AgentVersion, node.LastActive, node.CreatedAt, node.UpdatedAt)
	return err
}

func (r *NodeRepository) UpdateNode(node *Node) error {
	_, err := r.db.Exec("UPDATE nodes SET name = $1, hostname = $2, status = $3, operating_system = $4, agent_version = $5, last_active = $6, updated_at = $7 WHERE id = $8",
		node.Name, node.Hostname, node.Status, node.OperatingSystem, node.AgentVersion, node.LastActive, node.UpdatedAt, node.ID)
	return err
}

func (r *NodeRepository) DeleteNode(nodeID string) error {
	_, err := r.db.Exec("DELETE FROM nodes WHERE id = $1", nodeID)
	return err
}

func (r *NodeRepository) GetNodeByID(nodeID string) (*Node, error) {
	row := r.db.QueryRow("SELECT id, organization_id, name, hostname, status, operating_system, agent_version, last_active, created_at, updated_at FROM nodes WHERE id = $1", nodeID)

	var node Node
	if err := row.Scan(&node.ID, &node.OrganizationID, &node.Name, &node.Hostname, &node.Status, &node.OperatingSystem, &node.AgentVersion, &node.LastActive, &node.CreatedAt, &node.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &node, nil
}

func (r *NodeRepository) GetNodesByStatus(status Status) ([]*Node, error) {
	rows, err := r.db.Query("SELECT id, organization_id, name, hostname, status, operating_system, agent_version, last_active, created_at, updated_at FROM nodes WHERE status = $1", status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nodes []*Node
	for rows.Next() {
		var node Node
		if err := rows.Scan(&node.ID, &node.OrganizationID, &node.Name, &node.Hostname, &node.Status, &node.OperatingSystem, &node.AgentVersion, &node.LastActive, &node.CreatedAt, &node.UpdatedAt); err != nil {
			return nil, err
		}
		nodes = append(nodes, &node)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *NodeRepository) UpdateNodeStatus(nodeID string, status Status, lastActive time.Time) error {
	_, err := r.db.Exec("UPDATE nodes SET status = $1, last_active = $2, updated_at = $3 WHERE id = $4", status, lastActive, time.Now(), nodeID)
	return err
}

func (r *NodeRepository) SetPairingToken(nodeID, pairingTokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec("UPDATE nodes SET pairing_token_hash = $1, pairing_token_expires_at = $2, updated_at = $3 WHERE id = $4",
		pairingTokenHash, expiresAt, time.Now(), nodeID)
	return err
}

func (r *NodeRepository) GetNodeByPairingTokenHash(pairingTokenHash string) (*Node, *time.Time, error) {
	row := r.db.QueryRow("SELECT id, organization_id, name, hostname, status, operating_system, agent_version, last_active, created_at, updated_at, pairing_token_expires_at FROM nodes WHERE pairing_token_hash = $1", pairingTokenHash)

	var node Node
	var expiresAt *time.Time
	if err := row.Scan(&node.ID, &node.OrganizationID, &node.Name, &node.Hostname, &node.Status, &node.OperatingSystem, &node.AgentVersion, &node.LastActive, &node.CreatedAt, &node.UpdatedAt, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	return &node, expiresAt, nil
}

// CompletePairing atomically consumes the pairing token and issues the
// long-lived agent token, so a token can only ever be used to register once.
func (r *NodeRepository) CompletePairing(nodeID, agentTokenHash, hostname, operatingSystem, agentVersion string, lastActive time.Time) error {
	_, err := r.db.Exec(`UPDATE nodes
		SET agent_token_hash = $1, pairing_token_hash = NULL, pairing_token_expires_at = NULL,
			hostname = $2, operating_system = $3, agent_version = $4, status = $5, last_active = $6, updated_at = $7
		WHERE id = $8`,
		agentTokenHash, hostname, operatingSystem, agentVersion, ONLINE, lastActive, time.Now(), nodeID)
	return err
}

func (r *NodeRepository) GetNodeByAgentTokenHash(agentTokenHash string) (*Node, error) {
	row := r.db.QueryRow("SELECT id, organization_id, name, hostname, status, operating_system, agent_version, last_active, created_at, updated_at FROM nodes WHERE agent_token_hash = $1", agentTokenHash)

	var node Node
	if err := row.Scan(&node.ID, &node.OrganizationID, &node.Name, &node.Hostname, &node.Status, &node.OperatingSystem, &node.AgentVersion, &node.LastActive, &node.CreatedAt, &node.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &node, nil
}
