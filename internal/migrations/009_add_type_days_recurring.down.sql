-- Drop constraint and function
ALTER TABLE availability_rules DROP CONSTRAINT IF EXISTS check_days_of_week_range;
DROP FUNCTION IF EXISTS validate_days_of_week(INTEGER[]);

-- Re-add old day_of_week column
ALTER TABLE availability_rules ADD COLUMN day_of_week INTEGER;

-- Migrate data back: Take the first element of the array
-- Warning: This loses data if multiple days were set
UPDATE availability_rules 
SET day_of_week = days_of_week[1]
WHERE days_of_week IS NOT NULL AND array_length(days_of_week, 1) > 0;

-- Drop new columns
ALTER TABLE availability_rules 
    DROP COLUMN IF EXISTS days_of_week,
    DROP COLUMN IF EXISTS is_recurring,
    DROP COLUMN IF EXISTS type;
