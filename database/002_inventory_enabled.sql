-- add new column to the table
ALTER TABLE inventory
  ADD COLUMN enabled bool NOT NULL DEFAULT TRUE
