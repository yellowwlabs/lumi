package organizations

import "database/sql"

type OrganizationRepository struct {
	db *sql.DB
}

func NewOrganizationRepository(db *sql.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) GetOrganizationMembers(organizationID string) ([]*OrganizationMembers, error) {
	rows, err := r.db.Query("SELECT organization_id, user_id, role, created_at, updated_at FROM organization_members WHERE organization_id = $1", organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*OrganizationMembers
	for rows.Next() {
		var member OrganizationMembers
		if err := rows.Scan(&member.OrganizationID, &member.UserID, &member.Role, &member.CreatedAt, &member.UpdatedAt); err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

func (r *OrganizationRepository) CreateOrganizationMember(member *OrganizationMembers) error {
	_, err := r.db.Exec("INSERT INTO organization_members (organization_id, user_id, role, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		member.OrganizationID, member.UserID, member.Role, member.CreatedAt, member.UpdatedAt)
	return err
}

func (r *OrganizationRepository) GetOrganizationMember(organizationID, userID string) (*OrganizationMembers, error) {
	row := r.db.QueryRow("SELECT organization_id, user_id, role, created_at, updated_at FROM organization_members WHERE organization_id = $1 AND user_id = $2", organizationID, userID)

	var member OrganizationMembers
	if err := row.Scan(&member.OrganizationID, &member.UserID, &member.Role, &member.CreatedAt, &member.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &member, nil
}
