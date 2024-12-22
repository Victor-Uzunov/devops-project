ALTER TABLE todos
DROP CONSTRAINT IF EXISTS check_start_date_before_due_date;

ALTER TABLE todos
DROP COLUMN IF EXISTS start_date;
