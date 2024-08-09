CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE moderation_status AS ENUM ('created', 'approved', 'declined', 'on moderation');
CREATE TYPE user_type AS ENUM ('client', 'moderator');

CREATE TABLE houses (
    id SERIAL PRIMARY KEY,
    address TEXT NOT NULL,
    year INT NOT NULL,
    developer TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE subscribers (
    house_id INT REFERENCES houses(id) ON DELETE CASCADE,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE flats (
    id SERIAL PRIMARY KEY,
    house_id INT NOT NULL REFERENCES houses(id) ON DELETE CASCADE,
    price INT NOT NULL,
    rooms INT NOT NULL,
    status moderation_status
);

CREATE INDEX IF NOT EXISTS house_id_flat_id ON flats (house_id, id);

CREATE TABLE users (
    uid UUID  DEFAULT gen_random_uuid() PRIMARY KEY,
    email TEXT UNIQUE, 
    password TEXT,
    type user_type
);

CREATE INDEX IF NOT EXISTS user_email ON users (email);