alter table events add column "index" integer NOT NULL default 0;
alter table events add constraint unique_event unique ("type", txhash, blocknumber, "index");