CREATE TABLE IF NOT EXISTS items (
 id UUID PRIMARY KEY,
 owner_id TEXT NOT NULL,
 name TEXT NOT NULL CHECK (length(name) BETWEEN 1 AND 120),
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_items_owner_created ON items (owner_id, created_at DESC);
