CREATE TABLE IF NOT EXISTS liquidations (
  id serial PRIMARY KEY,
  eventId integer references events NOT NULL,
  borrowerId integer references accounts NOT NULL,
  liquidatorId integer references accounts NOT NULL,
  collateralTokenId integer references tokens NOT NULL,
  debtTokenId integer references tokens NOT NULL,
  amountRepayed bigint NOT NULL,
  amountSeized bigint NOT NULL
);
CREATE TABLE IF NOT EXISTS repays (
  id serial PRIMARY KEY,
  eventId integer references events NOT NULL,
  tokenId integer references tokens NOT NULL,
  borrowerId integer references accounts NOT NULL,
  amountRepayed bigint NOT NULL
);
CREATE TABLE IF NOT EXISTS borrows (
  id serial PRIMARY KEY,
  eventId integer references events NOT NULL,
  tokenId integer references tokens NOT NULL,
  borrowerId integer references accounts NOT NULL,
  amountBorrowed bigint NOT NULL
);