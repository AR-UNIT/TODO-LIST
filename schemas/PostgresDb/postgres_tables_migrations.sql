CREATE SCHEMA IF NOT EXISTS TODO;

CREATE TABLE TODO.clients (
                         id SERIAL PRIMARY KEY,                -- Auto-incrementing unique identifier
                         name VARCHAR(255) NOT NULL,           -- Name of the client
                         client_id VARCHAR(255) UNIQUE NOT NULL, -- Unique client identifier
                         client_secret VARCHAR(255) NOT NULL, -- Client secret (hashed for security)
                         created_at TIMESTAMP DEFAULT NOW()   -- Timestamp of creation
);


CREATE TABLE TODO.tasks (
                       id SERIAL PRIMARY KEY,
                       description TEXT NOT NULL,
                       completed BOOLEAN DEFAULT FALSE
);


INSERT INTO TODO.clients (name, client_id, client_secret, created_at)
VALUES
    ('admin', 'admin', 'password', '2024-12-30 21:22:54.859');