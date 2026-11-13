-- Revert availability_rules id column from UUID to SERIAL
-- Warning: This will result in new integer IDs being assigned.

-- Add new SERIAL column
ALTER TABLE availability_rules ADD COLUMN new_id SERIAL;

-- Drop primary key constraint
ALTER TABLE availability_rules DROP CONSTRAINT availability_rules_pkey;

-- Drop UUID id column
ALTER TABLE availability_rules DROP COLUMN id;

-- Rename new_id to id
ALTER TABLE availability_rules RENAME COLUMN new_id TO id;

-- Add primary key constraint on new SERIAL column
ALTER TABLE availability_rules ADD PRIMARY KEY (id);
