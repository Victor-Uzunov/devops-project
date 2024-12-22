alter table todos
    alter column due_date type timestamp using due_date::timestamp;

alter table todos
    alter column start_date type timestamp using start_date::timestamp;

