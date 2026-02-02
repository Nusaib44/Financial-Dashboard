CREATE TABLE IF NOT EXISTS daily_revenues (
  id UUID PRIMARY KEY,
  agency_id UUID NOT NULL REFERENCES agencies(id),
  date DATE NOT NULL,
  amount NUMERIC NOT NULL,
  source TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS daily_costs (
  id UUID PRIMARY KEY,
  agency_id UUID NOT NULL REFERENCES agencies(id),
  date DATE NOT NULL,
  amount NUMERIC NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('fixed', 'variable')),
  label TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT now()
);
