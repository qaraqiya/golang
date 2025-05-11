CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    category TEXT NOT NULL,
    amount FLOAT NOT NULL,
    description TEXT,
    date DATE NOT NULL
);
