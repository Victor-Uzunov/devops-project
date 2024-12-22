BEGIN;

ALTER TABLE users
    DROP COLUMN refresh_token,
    DROP COLUMN refresh_token_expiration;

COMMIT;
