
contract('MetaCoin', (accounts) => {

  let metacoinContract;

  describe('CONTRACT DEPLOYMENT', () => {
    it('Deploy Contribution contracts', (done) => {
      MetaCoin.new().then((result) => {
        metacoinContract = result;
        console.log("MetaCoin address:", metacoinContract.address)
        return metacoinContract.getBalance.call(accounts[0]);
      }).then((result) => {
        //metacoinContract = MelonToken.at(result);
        assert.equal(result, 1000);
        done();
      });
    });

    it('Send & check balances', (done) => {
      var account_one = accounts[0];
      var account_two = accounts[1];
      var sendamount = 10;

      return metacoinContract.sendCoin(account_two, sendamount).then((result) => {
        return metacoinContract.getBalance.call(account_one);
      }).then((result) => {
          assert.equal(result, 1000-sendamount);
          return metacoinContract.getBalance.call(account_two);
      }).then((result) => {
        assert.equal(result, sendamount);
        console.log("account_two balance:", result)
        done();
      });
    });
  });
});
