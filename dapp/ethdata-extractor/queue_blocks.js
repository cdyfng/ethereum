var Web3 = require("web3");
//var BigNumber = require("bignumber.js");
var winston = require("winston");

var pg = require("pg"),
  conn_string = "postgres://postgres@localhost:5432/blockchain_ethereum",
  db = new pg.Client(conn_string);
db.connect();

var RPC_SERVER = "http://localhost:8545";

var web3 = new Web3(new Web3.providers.HttpProvider(RPC_SERVER));

var BLOCK_REWARD = 5,
  FIRST_BLOCK = 3110000,
  MAXIMUM_BLOCK = web3.eth.blockNumber,
  BLOCKS_AT_A_TIME = 200,
  BLOCK_QUEUE_TIMEOUT = 1000;
var ONEGROUP_CNT = 0;


if (!web3.isConnected()) {
  winston.log("error", "web3 is not connected to the RPC");
  process.exit(1);
}

var add_address = function(address, is_contract) {
  //winston.log("info", "address: " + address)
  if (is_contract) {
    is_contract = 1;
  } else {
    is_contract = 0;
  }
  /*
  CREATE TABLE Addresses (address char(42),address_type integer,UNIQUE(address));
*/
  db.query("INSERT INTO Addresses (address, address_type) VALUES ($1, $2) ON CONFLICT DO NOTHING", [address, is_contract]);
};

var process_result = function(result) {
  add_address(result.miner);
  winston.log("info", [result.number, result.hash, result.timestamp, result.parentHash, result.nonce, result.miner, parseInt(result.difficulty), result.size, BLOCK_REWARD])
  //CREATE TABLE Blocks (block_number bigint, block_hash char(66),timestamp_utc timestamp, parent_hash char(66), nonce char(18), miner_addr char(42), difficulty bigint, size_bytes integer, block_reward integer, UNIQUE(block_number));
  //db.query("INSERT INTO Blocks (block_number, block_hash, timestamp_utc, parent_hash, nonce, miner_addr, difficulty, size_bytes, block_reward) VALUES ($1, $2, to_timestamp($3), $4, $5, $6, $7, $8, $9) ON CONFLICT DO NOTHING",
  db.query("INSERT INTO Blocks (block_number, block_hash, timestamp_utc, parent_hash, nonce, miner_addr, difficulty, size_bytes, block_reward) VALUES ($1, $2, to_timestamp($3), $4, $5, $6, $7, $8, $9)",
    [result.number, result.hash, result.timestamp, result.parentHash, result.nonce, result.miner, parseInt(result.difficulty), result.size, BLOCK_REWARD]);

  var cur_tx_type,
    cur_tx;
  for (var i = 0; i < result.transactions.length; i++) {
    cur_tx = result.transactions[i];
    if (cur_tx.to == null) {
      cur_tx_type = 2; // Contract Creation
    } else {
      /*if(web3.eth.getCode(cur_tx.to) == "0x"){
      	cur_tx_type = 0; // Person to Person
      } else {
      	cur_tx_type = 1; // Person to Contract
      }*/
      cur_tx_type = 0;
      add_address(cur_tx.to, cur_tx_type);
    }

    add_address(cur_tx.from);

    //winston.log("info", [cur_tx.hash, cur_tx.transactionIndex, cur_tx.input, cur_tx_type])
    //CREATE TABLE Transactions (tx_hash char(66), tx_index integer, extra_data text, transaction_type integer, UNIQUE(tx_hash));
    db.query("INSERT INTO Transactions (tx_hash, tx_index, extra_data, transaction_type) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
      [cur_tx.hash, cur_tx.transactionIndex, cur_tx.input, cur_tx_type]);

  }

}

var queue_block = function(block_number) {
    web3.eth.getBlock(block_number, true, function(error, result) {
      if (error) {
        winston.log("error", "Error getting block number: " + block_number + ". '" + error + "'");
      } else {
        if (result) {
          //winston.log("info", "Queueing Block #" + block_number + " (with " + result.transactions.length + " transactions).");
          //connection.publish("ethblocks", result);
          process_result(result);
        } else {
          winston.log("warn", "No block seen for #" + block_number);
        }
      }
      ONEGROUP_CNT--;
    });
  },
  queue_n_blocks = function(starting_block, n) {
    winston.log("info", "Queueing " + n + " blocks, from " + starting_block);
    for (var i = starting_block; i < (starting_block + n); i++) {
      if (i <= MAXIMUM_BLOCK) {
        ONEGROUP_CNT++;
        queue_block(i);
      } else {
        winston.log("info", "skipped #" + i);
      }
    }
  };

var process_new_block = function(error, block_hash) {
  if (error) {
    winston.log("error", "Error processing new block.");
  } else {
    web3.eth.getBlock(block_hash, true, function(error, result) {
      if (error) {
        winston.log("error", "Error getting block hash: " + block_hash + ". '" + error + "'");
      } else {
        if (result) {
          winston.log("info", "Queueing Block #" + result.number + " (with " + result.transactions.length + " transactions).");
          process_result(result);
        } else {
          winston.log("warn", "No block seen for #" + block_hash);
        }
      }
    });
  }
};

var repeat_queue_blocks = function(current_block) {
  var next_start_block = current_block;
  MAXIMUM_BLOCK = web3.eth.blockNumber;
  if (current_block < MAXIMUM_BLOCK) {
    if (ONEGROUP_CNT <= 0) {
      winston.log("info", "finished , now restart " + ONEGROUP_CNT)
      queue_n_blocks(current_block, BLOCKS_AT_A_TIME);
      next_start_block += BLOCKS_AT_A_TIME + 1;
    } else {
      winston.log("info", "skip  and wait " + ONEGROUP_CNT)
    }

    setTimeout(repeat_queue_blocks.bind(this, next_start_block), BLOCK_QUEUE_TIMEOUT)
  //setTimeout(repeat_queue_blocks.bind(this, (current_block+BLOCKS_AT_A_TIME+1)), BLOCK_QUEUE_TIMEOUT)
  } else {
    winston.log("info", "We have caught up to the latest block #" + current_block + ", now listening...");
    web3.eth.filter("latest").watch(process_new_block);
  }
};
setTimeout(repeat_queue_blocks.bind(this, FIRST_BLOCK), BLOCK_QUEUE_TIMEOUT);
