BEGIN;

DROP TRIGGER IF EXISTS update_lists_timestamp ON lists;
DROP TRIGGER IF EXISTS update_todos_timestamp ON todos;
DROP TRIGGER IF EXISTS update_users_timestamp ON users;

DROP FUNCTION IF EXISTS update_timestamp();

COMMIT;