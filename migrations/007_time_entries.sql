CREATE TABLE IF NOT EXISTS time_entries (
  id UUID PRIMARY KEY,
  agency_id UUID NOT NULL REFERENCES agencies(id),
  client_id UUID NULL REFERENCES clients(id),
  date DATE NOT NULL,
  hours NUMERIC NOT NULL CHECK (hours > 0),
  created_at TIMESTAMP DEFAULT now()
);
