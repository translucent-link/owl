CREATE TABLE IF NOT EXISTS events (
   id serial PRIMARY KEY,
   type VARCHAR (50) NOT NULL,
   txhash VARCHAR (255) NOT NULL,
   blocknumber integer NOT NULL,
   occuredAt timestamp NOT NULL,
   protocolInstanceId integer references protocol_instances,
   eventDefinitionId integer references event_definitions
);
CREATE INDEX IF NOT EXISTS events_type ON events (type);