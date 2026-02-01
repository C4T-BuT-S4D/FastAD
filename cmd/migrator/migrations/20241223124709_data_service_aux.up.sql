SET statement_timeout = 0;

--bun:split

ALTER TABLE game_state
    ADD CONSTRAINT hardness_non_zero CHECK (hardness > 0);
