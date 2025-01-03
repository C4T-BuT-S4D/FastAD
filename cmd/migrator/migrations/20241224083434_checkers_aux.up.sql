SET statement_timeout = 0;

--bun:split

-- For faster flag picking for GET.
CREATE INDEX IF NOT EXISTS idx_flags_team_service_round ON flags (team_id, service_id, round);

--bun:split

--For faster slac processing
CREATE INDEX IF NOT EXISTS idx_checker_executions_created_at ON checker_executions (created_at);

--bun:split

--For faster attack_data handle
CREATE INDEX IF NOT EXISTS idx_attack_data_snapshots ON attack_data_snapshots (created_at DESC);
