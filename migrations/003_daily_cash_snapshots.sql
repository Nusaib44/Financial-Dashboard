CREATE TABLE IF NOT EXISTS daily_cash_snapshots (
  id UUID PRIMARY KEY,
  agency_id UUID NOT NULL REFERENCES agencies(id),
  date DATE NOT NULL,
  cash_balance NUMERIC NOT NULL,
  created_at TIMESTAMP DEFAULT now(),
  UNIQUE (agency_id, date)
);
