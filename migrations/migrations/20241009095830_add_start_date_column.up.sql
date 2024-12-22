BEGIN;

ALTER TABLE todos
ADD COLUMN start_date TIMESTAMP;

ALTER TABLE todos
ADD CONSTRAINT check_start_date_before_due_date
CHECK ( start_date IS NULL OR due_date IS NULL OR start_date <= due_date );

COMMIT;