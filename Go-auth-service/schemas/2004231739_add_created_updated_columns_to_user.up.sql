ALTER TABLE account ADD COLUMN createdAt TIMESTAMP;
ALTER TABLE account ADD COLUMN updatedAt TIMESTAMP;

ALTER TABLE account ALTER COLUMN createdAt SET DEFAULT current_timestamp;
ALTER TABLE account ALTER COLUMN updatedAt SET DEFAULT current_timestamp;