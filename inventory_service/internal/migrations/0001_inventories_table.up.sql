-- Create table Inventory with product_id as primary key and amount as integer not null default 0
CREATE TABLE inventories
(
    product_id VARCHAR PRIMARY KEY,
    amount     INTEGER   NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);