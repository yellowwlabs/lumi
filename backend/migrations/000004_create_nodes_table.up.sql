CREATE TABLE IF NOT EXISTS nodes (
    id                          UUID PRIMARY KEY,
    organization_id             UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name                        TEXT NOT NULL,
    hostname                    TEXT NOT NULL DEFAULT '',
    status                      TEXT NOT NULL DEFAULT 'pending',
    operating_system            TEXT NOT NULL DEFAULT '',
    agent_version               TEXT NOT NULL DEFAULT '',
    last_active                 TIMESTAMPTZ NOT NULL DEFAULT now(),
    pairing_token_hash          TEXT,
    pairing_token_expires_at    TIMESTAMPTZ,
    agent_token_hash            TEXT,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_nodes_organization_id ON nodes(organization_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_nodes_pairing_token_hash ON nodes(pairing_token_hash) WHERE pairing_token_hash IS NOT NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_nodes_agent_token_hash ON nodes(agent_token_hash) WHERE agent_token_hash IS NOT NULL;
