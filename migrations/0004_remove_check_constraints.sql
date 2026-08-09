-- Remove hardcoded CHECK constraints to allow dynamic values from master tables

-- Drop constraints on activity_logs
ALTER TABLE activity_logs DROP CONSTRAINT IF EXISTS activity_logs_pic_check;
ALTER TABLE activity_logs DROP CONSTRAINT IF EXISTS activity_logs_status_check;
ALTER TABLE activity_logs DROP CONSTRAINT IF EXISTS activity_logs_category_check;

-- Also fix pic column length (was too short for actual names)
ALTER TABLE activity_logs ALTER COLUMN pic TYPE VARCHAR(100);
ALTER TABLE activity_logs ALTER COLUMN status TYPE VARCHAR(100);
ALTER TABLE activity_logs ALTER COLUMN category TYPE VARCHAR(100);
