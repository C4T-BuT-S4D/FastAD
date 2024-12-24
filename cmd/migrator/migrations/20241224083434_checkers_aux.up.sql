SET statement_timeout = 0;

--bun:split

-- For faster flag picking for GET.
CREATE INDEX IF NOT EXISTS idx_flags_team_service_round ON flags (team_id, service_id, round);

--bun:split

--For faster scoreboard processing
CREATE INDEX IF NOT EXISTS idx_checker_executions_created_at ON checker_executions (created_at);
