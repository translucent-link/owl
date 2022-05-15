CREATE TABLE IF NOT EXISTS event_definitions (
   id serial PRIMARY KEY,
   topicName VARCHAR (50) NOT NULL,
   protocolId integer references protocols,
   topicHashHex VARCHAR (255) NOT NULL,
   abiSignature TEXT
);