CREATE TABLE IF NOT EXISTS cash_snapshots (
    id UUID PRIMARY KEY,
    agency_id UUID NOT NULL,
    date DATE NOT NULL,
    cash_balance NUMERIC NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (agency_id, date)
);
