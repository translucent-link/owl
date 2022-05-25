CREATE TABLE IF NOT EXISTS borrows (
  id serial PRIMARY KEY,
  eventId integer references events NOT NULL,
  tokenId integer references tokens NOT NULL,
  borrowerId integer references accounts NOT NULL,
  amountBorrowed bigint NOT NULL
);