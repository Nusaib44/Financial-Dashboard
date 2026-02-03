-- Add category to daily_costs
ALTER TABLE daily_costs ADD COLUMN IF NOT EXISTS category TEXT NOT NULL DEFAULT 'other' CHECK (category IN ('people', 'tools', 'other'));
