Etherum data extractor
========

- extract data from the Ethereum blockchains
- save data in postgresql

## Step 1

install postgresql

create database blockchain_ethereum;

\c blockchain_ethereum

CREATE TABLE Addresses (address char(42),address_type integer,UNIQUE(address));

CREATE TABLE Blocks (block_number bigint, block_hash char(66),timestamp_utc timestamp, parent_hash char(66), nonce char(18), miner_addr char(42), difficulty bigint, size_bytes integer, block_reward integer, UNIQUE(block_number));

CREATE TABLE Transactions (tx_hash char(66), tx_index integer, extra_data text, transaction_type integer, UNIQUE(tx_hash));
## Step 2

cd ethdata-extractor

npm install

## Step 3

run an Ethereum node  and make sure rpc is available

## Step 4

node queue_blocks.js
