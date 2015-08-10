CREATE TABLE IF NOT EXISTS players (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  team VARCHAR(255) NOT NULL,
  vs_team VARCHAR(255) NOT NULL,
  pos VARCHAR(255) NOT NULL,
  points FLOAT,
  price INTEGER
);

CREATE TABLE IF NOT EXISTS teams (
  id SERIAL PRIMARY KEY,
  pitcher INTEGER NOT NULL,
  catcher INTEGER NOT NULL,
  first INTEGER NOT NULL,
  second INTEGER NOT NULL,
  third INTEGER NOT NULL,
  short INTEGER NOT NULL,
  of1 INTEGER NOT NULL,
  of2 INTEGER NOT NULL,
  of3 INTEGER NOT NULL,
  points FLOAT 
);