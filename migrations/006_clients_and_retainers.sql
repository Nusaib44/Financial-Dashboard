CREATE TABLE IF NOT EXISTS clients (
  id UUID PRIMARY KEY,
  agency_id UUID NOT NULL REFERENCES agencies(id),
  name TEXT NOT NULL,
  status TEXT NOT NULL CHECK (status IN ('active', 'paused', 'ended')),
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS retainers (
  id UUID PRIMARY KEY,
  agency_id UUID NOT NULL REFERENCES agencies(id),
  client_id UUID NOT NULL REFERENCES clients(id),
  monthly_amount NUMERIC NOT NULL,
  active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT now()
);
