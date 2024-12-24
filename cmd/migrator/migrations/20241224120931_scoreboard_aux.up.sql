SET statement_timeout = 0;

--bun:split

-- For faster except-filtering.
CREATE INDEX IF NOT EXISTS idx_processor_states_processor_id_entity_id ON processor_states (processor_id, entity_id);

--bun:split

-- For faster restore.
CREATE INDEX IF NOT EXISTS idx_processor_states_processr_id_id ON processor_states (processor_id, id);
