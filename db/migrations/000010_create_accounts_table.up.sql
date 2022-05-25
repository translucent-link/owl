CREATE TABLE IF NOT EXISTS accounts (
   id serial PRIMARY KEY,
   address VARCHAR (255) NOT NULL
);
ALTER TABLE accounts ADD CONSTRAINT uniq_accounts_address UNIQUE (address);