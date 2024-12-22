BEGIN;

CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_lists_timestamp
    BEFORE UPDATE ON lists
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE OR REPLACE TRIGGER update_todos_timestamp
    BEFORE UPDATE ON todos
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE OR REPLACE TRIGGER update_users_timestamp
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

COMMIT;

