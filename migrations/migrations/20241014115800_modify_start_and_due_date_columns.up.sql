BEGIN;

alter table todos
    alter column due_date type timestamptz using due_date::timestamptz;

alter table todos
    alter column start_date type timestamptz using start_date::timestamptz;

COMMIT;