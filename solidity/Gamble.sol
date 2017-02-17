pragma solidity ^0.4.6;

//简单对赌协议
//打以太币到合约，满足两个账户，谁hash大，谁得到所有以太币

contract Gamble{
  bool public inSupport;
  uint256 public totalEther;
  mapping(address => uint256) public balanceOf;
  event FundTransfer(address backer, uint amount, bool isContribution);

  function () payable {
      uint amount = msg.value;
      balanceOf[msg.sender] = amount;
      totalEther += amount;
      FundTransfer(msg.sender, amount, true);
      //checkGame();
      if(inSupport){
          //返回所有ether给  出发此处代码到sender
          FundTransfer(msg.sender, totalEther, false);
      }
      inSupport = !inSupport;
  }

  function checkGame() {

    if(inSupport == true){
        //返回所有ether给  出发此处代码到sender
        FundTransfer(msg.sender, totalEther, false);
    }
  }

  function safeWithdrawal() {
    uint amount = balanceOf[msg.sender];
    balanceOf[msg.sender] = 0;
    if (amount > 0) {
        if (msg.sender.send(totalEther)) {
            FundTransfer(msg.sender, totalEther, false);
        } else {
            balanceOf[msg.sender] = amount;
        }
    }
  }

  function finalBetCheck() {

  }


}
