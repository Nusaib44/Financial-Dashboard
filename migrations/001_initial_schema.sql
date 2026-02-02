CREATE TABLE IF NOT EXISTS founders (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS agencies (
    id UUID PRIMARY KEY,
    owner_user_id UUID NOT NULL,
    name TEXT NOT NULL,
    base_currency TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
