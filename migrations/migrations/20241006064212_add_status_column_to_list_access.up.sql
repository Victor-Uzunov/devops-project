BEGIN;

CREATE TYPE access_status AS ENUM ('owner', 'accepted', 'pending');

ALTER TABLE list_access
    ADD COLUMN status access_status DEFAULT 'pending';

UPDATE list_access
SET status = 'pending'
WHERE access_level IN ('reader', 'writer', 'admin');

ALTER TABLE list_access
    ALTER COLUMN status SET NOT NULL;

COMMIT;