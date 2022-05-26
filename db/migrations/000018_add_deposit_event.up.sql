alter table events 
add column depositTokenId integer NULL references tokens,
add column amountDeposited numeric NOT NULL default 0,
add column depositorAccountId integer NULL references accounts;