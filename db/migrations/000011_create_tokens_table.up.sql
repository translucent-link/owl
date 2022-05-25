CREATE TABLE IF NOT EXISTS tokens (
   id serial PRIMARY KEY,
   address VARCHAR (255) NOT NULL,
   name VARCHAR(50),
   ticker VARCHAR(10),
   chainId integer references chains
);
ALTER TABLE tokens ADD CONSTRAINT uniq_token_chain UNIQUE (address, chainId);