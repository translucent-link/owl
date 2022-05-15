CREATE TABLE IF NOT EXISTS protocol_instances (
   id serial PRIMARY KEY,
   contractAddress VARCHAR (50) NOT NULL,
   protocolId integer references protocols,
   chainId integer references chains
);