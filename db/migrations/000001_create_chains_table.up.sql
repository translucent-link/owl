CREATE TABLE IF NOT EXISTS chains(
   id serial PRIMARY KEY,
   name VARCHAR (50) UNIQUE NOT NULL,
   rpcUrl VARCHAR (255) NOT NULL,
   blockFetchSize integer NOT NULL default 10000
);