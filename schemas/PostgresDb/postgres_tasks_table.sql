CREATE TABLE tasks (
                       id SERIAL PRIMARY KEY,
                       description TEXT NOT NULL,
                       completed BOOLEAN DEFAULT FALSE
);
