CREATE TABLE clients (
                         id SERIAL PRIMARY KEY,                -- Auto-incrementing unique identifier
                         name VARCHAR(255) NOT NULL,           -- Name of the client
                         client_id VARCHAR(255) UNIQUE NOT NULL, -- Unique client identifier
                         client_secret VARCHAR(255) NOT NULL, -- Client secret (hashed for security)
                         created_at TIMESTAMP DEFAULT NOW()   -- Timestamp of creation
);
