SET statement_timeout = 0;

--bun:split

-- For faster flag picking for GET.
CREATE INDEX IF NOT EXISTS idx_flags_team_service_round ON flags (team_id, service_id, round);
