BEGIN;

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    email VARCHAR(255)  UNIQUE NOT NULL,
    github_id VARCHAR(255) UNIQUE NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('reader', 'writer', 'admin')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE lists (
    id UUID PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    visibility VARCHAR(50) NOT NULL CHECK (visibility IN ('public', 'private', 'shared')) DEFAULT 'private',
    tags jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE todos (
    id UUID PRIMARY KEY NOT NULL,
    list_id UUID REFERENCES lists(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    due_date TIMESTAMP,
    priority VARCHAR(50) CHECK (priority IN ('low', 'medium', 'high')),
    tags jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    assigned_to UUID
);

CREATE TABLE list_access (
    list_id UUID REFERENCES lists(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    access_level VARCHAR(50) NOT NULL CHECK (access_level IN ('reader', 'writer', 'admin')) DEFAULT 'reader',
    PRIMARY KEY (list_id, user_id)
);

CREATE INDEX idx_lists_owner_id ON lists(owner_id);
CREATE INDEX idx_todos_list_id ON todos(list_id);

COMMIT;