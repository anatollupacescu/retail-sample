-- add new column to the table
ALTER TABLE recipe
  ADD COLUMN enabled bool NOT NULL DEFAULT TRUE
