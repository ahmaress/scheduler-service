-- Re-add the unique constraint
-- Note: This might fail if there are duplicate records. Manual cleanup would be required.
ALTER TABLE availability_rules
ADD CONSTRAINT uniq_user_day UNIQUE (user_id, day_of_week);
