-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- This is test table. Remove this table and replace with your own tables. 

-- Enable the pgcrypto extension for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create the estate table
CREATE TABLE estate (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    length INTEGER NOT NULL,
    width INTEGER NOT NULL
);

-- Create the tree table
CREATE TABLE tree (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    estate_id UUID NOT NULL REFERENCES estate(id),
    x_axis INTEGER NOT NULL,
    y_axis INTEGER NOT NULL,
    height INTEGER NOT NULL
);