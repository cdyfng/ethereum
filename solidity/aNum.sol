pragma solidity ^0.4.0;

contract Contract {
  uint public num;
  event aNum(uint x);
  function Contract(uint x) {
    aNum(x);
    num = x;
  }
}
