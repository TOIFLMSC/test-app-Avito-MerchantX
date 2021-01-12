CREATE TABLE offers (
    id bigserial PRIMARY KEY,
    seller VARCHAR NOT NULL,
    offer_id VARCHAR,
    name VARCHAR,
    price bigint,
    quantity bigint,
    available VARCHAR
);