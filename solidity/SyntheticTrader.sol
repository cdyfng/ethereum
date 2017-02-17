pragma solidity ^0.4.6;
contract SyntheticTrader {

    // Work in process - EEEERRRROOORRRSSS

    // 'Simple' list of open ordåers
    // Sorted by price in buy and sell orders
    // If you sell more than you have - a security is placed
    // The security is released as soon you are back to 0

    // Heavily used links
    // https://chriseth.github.io/browser-solidity/
    // http://solidity.readthedocs.org/en/latest/

    int public No_Accounts;    // Max number of accounts
    int public No_Sell_Orders; // Max number of sell orders
    int public No_Buy_Orders;  // Max number of buy orders
    int public Amount;         // Amount of stock to be traded
    int public Price;          // Price of the stock to be traded
    int public Ref_Price;      // Refrence price for sell order securety
    int public Ref_Price_Block_No; // Change the reference price only every 100 bloks

    int public sU = 10**18;    // 1 Unit is 1/1e18 of a shear (sU = smallest Unit)

    mapping (address => int) public Own_Account_No; // Account number
    mapping (address => int) public Own_Funds;      // Funds of the trader in Wei (access by trader)
    mapping (address => int) public Own_Security;   // Security of the trader in Wei (no access by trader)
    mapping (address => int) public Own_Amount;     // Amount on Stock in Stock/pU
    mapping (address => int) public Own_Amount_Sell_Order;// Amount on Stock in sell orders in Stock/pU
    mapping (address => int) public Own_FeedBack;   // For debugging

    struct Account
    {
        address Account_Address;
    }
    mapping (int => Account) public Accounts;

    struct Sell
    {
       int Amount;
       int Price;
       address Address;
       int Security;
    }
    mapping (int => Sell) public Sells;

    struct Buy
    {
       int Amount;
       int Price;
       address Address;
    }
    mapping (int => Buy) public Buys;

    function SyntheticTrader() {
       // Initialization
       No_Sell_Orders = 0;                  // Start without orders
       No_Buy_Orders  = 0;
    }

    function () payable {// Send Ether to the contract

      Own_Funds[msg.sender] += int(msg.value); // Add Funds in Wei

      if (Own_Account_No[msg.sender] == 0){
        No_Accounts += 1;
        Own_Account_No[msg.sender] = No_Accounts;
        Accounts[No_Accounts].Account_Address = msg.sender;
      }

    }

    function Sell_Order(int Amount_in_sU, int Price_in_Wei) { // Sell order

        Amount = Amount_in_sU;
        Price  = Price_in_Wei;

        Own_FeedBack[msg.sender] =11; // 11 = Sell_Order

        while (Amount > 0 && Price > 0){
            if (Own_Amount[msg.sender] > 0 || Own_Funds[msg.sender] > 0 ){

                if (Buys[No_Buy_Orders].Price >= Price && No_Buy_Orders > 0) { // Sell if price is higher than ask
                    // Sell

                    Sell_from_List();

                } else {// to low in price
                    // Create Sell order with the rest Amount

                    Create_Sell_Order();

                }
            } else {
                Amount = 0;
            }
        }

    }

    function Buy_Order(int Amount_in_sU, int Price_in_Wei) { // New Buy order
        Amount = Amount_in_sU;
        Price  = Price_in_Wei;
        Buy_Order_(msg.sender,0); // 0 = no margin call
    }

    function Buy_Order_(address msg_sender, int Margin_Call) internal { // New Buy order

        Own_FeedBack[msg.sender] = 12; // 12 = Buy_Order

        while (Amount > 0 && Price > 0){
            if (Own_Funds[msg_sender] + Own_Security[msg_sender] > 0){
                if (Sells[No_Sell_Orders].Price <= Price && No_Sell_Orders > 0) { // Buy if price is lower than ask
                    // Buy

                    Buy_from_List(msg_sender);

                } else {// to high in price
                    if (Margin_Call == 0){
                    // Create Sell order with the rest Amount
                        Create_Buy_Order(msg_sender);
                    } else {
                    // Margin Call
                        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 15; // 15 = Margin Call
                        Own_Funds[msg_sender]-=Amount * Price / sU;
                        Own_Funds[0000000000]+=Amount * Price / sU;
                        Own_Amount[msg_sender] =0;
                        Own_Amount[0000000000]-=Amount;
                        Create_Buy_Order(0000000000);
                    }
                }
            } else {
                Amount = 0;
            }
        }

    }

    function Cancel_Order(int Price_in_Wei) { // Cancle all orders

        Own_FeedBack[msg.sender] = 13; // 13 = Cancel_Order

        Price  = Price_in_Wei;

        if (Price <= Buys[No_Buy_Orders].Price && No_Buy_Orders > 0){

            Cancel_Buy_Order();

        }

        if (Price >= Sells[No_Sell_Orders].Price && No_Sell_Orders > 0){

            Cancel_Sell_Order();

        }

    }

    function Margin_Call(){
        // check if there is an account with not enough security
        address addr;
        int Margin;
        int Own_Amount_Debt;

        for (int i = No_Accounts; i>0; i--){
            addr = Accounts[i].Account_Address;
            Own_Amount_Debt = min(Own_Amount[addr] + Own_Amount_Sell_Order[addr], 0) * (-1);
            if (Own_Amount_Debt > 0 && Sells[No_Sell_Orders].Price > 0){
                Margin = Own_Security[addr] * sU / Own_Amount_Debt;
                if (Sells[No_Sell_Orders].Price > Margin * 9 / 10){
                    Own_FeedBack[addr] = 15; // 15 = Margin Call
                    // Caller gets some of the security for the call
                    Own_Funds[msg.sender] += Own_Security[addr] * max(sU - Sells[No_Sell_Orders].Price * sU / Margin, 0) / sU / 2;
                    Own_Security[addr]    -= Own_Security[addr] * max(sU - Sells[No_Sell_Orders].Price * sU / Margin, 0) / sU / 2;
                    // Margin Call
                    Amount = Own_Amount_Debt;
                    Price  = Own_Security[addr] * sU / Own_Amount_Debt;
                    Own_Funds[addr] = Own_Funds[addr] + Own_Security[addr];
                    Own_Security[addr] = 0;
                    Buy_Order_(addr, 1); // Address and 1 for Margin Call

                    i = 0 ; // Exit
                }
            }
        }
    }

    function Withdraw_All_Funds_With_Error() { // Withdraw all the free funds of the trader

        Own_FeedBack[msg.sender] =14; // 14 = Withdraw_All_Funds

        if (Own_Funds[msg.sender]>0){
            //msg.sender.send(uint(Own_Funds[msg.sender]));
            //msg.sender.send(uint256(Own_Funds[msg.sender]));
            Own_Funds[msg.sender]=0;
            Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 01; // 1401 = Send funds
        }
    }

// ------------------------------------------------------------------------------
// internal routines
// ------------------------------------------------------------------------------

// ------------------------------------------------------------------------------
// Buy / Sell from List
// ------------------------------------------------------------------------------

    function Sell_from_List() internal { //  internal

        int List_Amount = Buys[No_Buy_Orders].Amount;
        int List_Price  = Buys[No_Buy_Orders].Price;
        int Max_Amount;

        int Transfer_Amount = min(Amount,List_Amount);
        int Own_Amount_Credit = max(Own_Amount[msg.sender],0);
        int Sell_Amount = min(Transfer_Amount, Own_Amount_Credit);

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 21; // 21 = Sell_from_List

        // First sell the Amount the trader has and get the funds
        Own_Funds[msg.sender] += Sell_Amount * List_Price / sU;

        // Then sell the remaining amount with funds as security

        int Pay_Amount = min(Transfer_Amount - Sell_Amount, Own_Funds[msg.sender]  * sU / List_Price);

        Own_Security[msg.sender] += Pay_Amount * List_Price / sU; // form buyer
        Own_Security[msg.sender] += Pay_Amount * List_Price / sU; // from seller
        Own_Funds[msg.sender]    -= Pay_Amount * List_Price / sU;

        Amount                   -= Pay_Amount + Sell_Amount;
        Own_Amount[msg.sender]   -= Pay_Amount + Sell_Amount;

        if (Amount <= 0){
            Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 90; // 2190 = Sell_from_List
            if (Pay_Amount > 0) {
             Own_FeedBack[msg.sender] += 1; // Sell with Funds 2191
            }
            if (Sell_Amount > 0) {
             Own_FeedBack[msg.sender] += 2; // Sell with Amount 2192  // 2193 Funds + Amount
            }
        }

        Sell_from_List_send_Buyer(Sell_Amount + Pay_Amount);
        Sell_from_List_edit_List(Sell_Amount + Pay_Amount);
        Set_Ref_Price(List_Price);
   }

    function Buy_from_List(address msg_sender) internal { //  internal

        int List_Amount = Sells[No_Sell_Orders].Amount;
        int List_Price  = Sells[No_Sell_Orders].Price;
        int Max_Amount;

        int Transfer_Amount = min(Amount,List_Amount);
        int Own_Amount_Debt = min(Own_Amount[msg.sender]+ Own_Amount_Sell_Order[msg.sender], 0) * (-1);

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 22; // 22 = Buy_from_List

        if (Own_Amount_Debt == 0) { // No Debt
            // trader buys only with funds // Own_Security = 0

            Max_Amount = Own_Funds[msg.sender] * sU / List_Price;

        } else { // Own_Amount < 0  Security > 0  // Debt

            if (List_Price <= Own_Security[msg.sender] * sU / Own_Amount_Debt){
                // trader can buy with the security
                Max_Amount = (Own_Security[msg.sender] + Own_Funds[msg.sender]) * sU / List_Price;
            }else{
                // trader has to add funds to relese his security
                Max_Amount = Own_Funds[msg.sender] * sU / (List_Price - Own_Security[msg.sender] * sU / Own_Amount_Debt);
                if (Own_Amount_Debt<Max_Amount){
                    int rem_Funds = Own_Funds[msg.sender] - (List_Price * Own_Amount_Debt / sU - Own_Security[msg.sender]);
                    Max_Amount = Own_Amount_Debt + rem_Funds * sU / List_Price;
                }
            }
        }

        if (Transfer_Amount >= Max_Amount){
            // exit because end of funds
            Transfer_Amount = Max_Amount;
            Amount = 0;
            Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 90; // 2290 = Buy_from_List - End of Funds
        }else{
            Amount = Amount - Transfer_Amount;
        }

        if (Own_Amount[msg.sender] < 0) {
            int Free_Security = Own_Security[msg.sender] * min(Transfer_Amount * sU / Own_Amount_Debt, sU) / sU;
            Own_Funds[msg.sender]    += max(Free_Security, 0);
            Own_Security[msg.sender] -= max(Free_Security, 0);
        }
        Own_Funds[msg.sender]  -=  Transfer_Amount * List_Price / sU;
        Own_Amount[msg.sender] +=  Transfer_Amount;

        Buy_from_List_send_Seller(Transfer_Amount);
        Buy_from_List_edit_List(Transfer_Amount);
        Set_Ref_Price(List_Price);
    }

    function Sell_from_List_send_Buyer(int Transfer_Amount) internal { // Internal

        address List_msg_sender = Buys[No_Buy_Orders].Address;

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 23; // 23 = Sell_from_List_send_Buyer

        if (Own_Amount[List_msg_sender] + Transfer_Amount > 0) {

            Own_Funds[List_msg_sender]   += Own_Security[List_msg_sender];
            Own_Security[List_msg_sender] = 0;

        } else { // List_Own_Amount + Transfer_Amount =< 0

            int Free_Security = Own_Security[List_msg_sender] * Transfer_Amount / (-Own_Amount[List_msg_sender]);
            Own_Funds[List_msg_sender]    += Free_Security;
            Own_Security[List_msg_sender] -= Free_Security;
        }

        Own_Amount[List_msg_sender] += Transfer_Amount;

    }

    function Buy_from_List_send_Seller(int Transfer_Amount) internal { // Internal

        address List_msg_sender = Sells[No_Sell_Orders].Address;

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 24; // 24 = Buy_from_List_send_Seller

        if (Own_Security[List_msg_sender] == 0) {// Seller has no dept

            Own_Funds[List_msg_sender]    += Sells[No_Sell_Orders].Price * Transfer_Amount / sU;

        } else {// Seller has dept

            Own_Security[List_msg_sender] += Sells[No_Sell_Orders].Price * Transfer_Amount / sU;

            if (Sells[No_Sell_Orders].Security > 0){
                Own_Security[List_msg_sender] += Sells[No_Sell_Orders].Security * Transfer_Amount / Sells[No_Sell_Orders].Amount;
                Own_Amount_Sell_Order[List_msg_sender] -= Transfer_Amount;
            }
        }

    }

    function Sell_from_List_edit_List(int Transfer_Amount) internal { // Internal

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 25; // 25 = Sell_from_List_edit_List

        if (Transfer_Amount < Buys[No_Buy_Orders].Amount) {
            // If some more available
            // Don't close the List entry
            Buys[No_Buy_Orders].Amount -= Transfer_Amount;
        } else {
            // Close the order
            Buys[No_Buy_Orders].Amount  = 0;
            Buys[No_Buy_Orders].Price   = 0;
            Buys[No_Buy_Orders].Address = 0;
            No_Buy_Orders -= 1;
        }
    }

    function Buy_from_List_edit_List(int Transfer_Amount) internal { // Internal

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 26; // 26 = Buy_from_List_edit_List

        if (Transfer_Amount < Sells[No_Sell_Orders].Amount) {
            // If some more is available
            // Don't close the list entry
            Sells[No_Sell_Orders].Security  = Sells[No_Sell_Orders].Security * (Sells[No_Sell_Orders].Amount - Transfer_Amount) / Sells[No_Sell_Orders].Amount;
            Sells[No_Sell_Orders].Amount   -= Transfer_Amount;
        } else {
            // Close the order
            Sells[No_Sell_Orders].Amount   = 0;
            Sells[No_Sell_Orders].Price    = 0;
            Sells[No_Sell_Orders].Address  = 0;
            Sells[No_Sell_Orders].Security = 0;
            No_Sell_Orders -= 1;
        }
    }

// ------------------------------------------------------------------------------
// Create Order
// ------------------------------------------------------------------------------

    function Create_Sell_Order() internal { // internal

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 31; // 31 = Create_Sell_Order

        int Security = 0;
        int Transfer_Amount = 0;
        int Own_Amount_Credit = max(Own_Amount[msg.sender],0);

        if (Own_Amount_Credit > 0) {
            Transfer_Amount = min(Amount, Own_Amount_Credit);
            Add_Sell_Order(Transfer_Amount, 0); // 0 = No Security
            Amount = Amount - Transfer_Amount;
        }
        if (Amount > 0) {
            Amount = min(Amount, Own_Funds[msg.sender] * sU / Ref_Price);

            Security = Amount * Ref_Price / sU;
            Own_Funds[msg.sender]    -= Security;

            Own_Amount[msg.sender] -= Amount;
            Own_Amount_Sell_Order[msg.sender] += Amount;

            Add_Sell_Order(Transfer_Amount, Security);

        }
        Amount = 0;
    }

    function Create_Buy_Order(address msg_sender) internal { // internal

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 32; // 32 = Create_Buy_Order

        Amount = min(Amount, Own_Funds[msg_sender] * sU / Price);

        if (Amount > 0) {

            Own_Funds[msg_sender] -= Amount * Price / sU;

            Add_Buy_Order(msg_sender);

        }
        Amount = 0;
    }

    function Add_Sell_Order(int Transfer_Amount, int Security) internal { // internal

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 33; // 33 = Add_Sell_Order

        No_Sell_Orders += 1;

        for (int i = No_Sell_Orders; i>0; i--){

            if (Sells[i-1].Price > Price || i == 1){

                Sells[i].Price      = Price;
                Sells[i].Amount     = Transfer_Amount;
                Sells[i].Address    = msg.sender;
                Sells[i].Security   = Security;

                Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + i; // Added Sell Order

                i = 0;

            }else{
                Sells[i].Price      = Sells[i-1].Price;
                Sells[i].Amount     = Sells[i-1].Amount;
                Sells[i].Address    = Sells[i-1].Address;
                Sells[i].Security   = Sells[i-1].Security;
            }
        }
    }

    function Add_Buy_Order(address msg_sender) internal { // internal

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 34; // 34 = Add_Buy_Order

        No_Buy_Orders += 1;

        for (int i = No_Buy_Orders; i>0; i--){

            if (Buys[i-1].Price < Price || i == 1){

                Buys[i].Price      = Price;
                Buys[i].Amount     = Amount;
                Buys[i].Address    = msg_sender;

                Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + i; // Added Buy Order

                i = 0;

            }else{
                Buys[i].Price      = Buys[i-1].Price;
                Buys[i].Amount     = Buys[i-1].Amount;
                Buys[i].Address    = Buys[i-1].Address;
            }
        }
    }

// ------------------------------------------------------------------------------
// Cancel Order
// ------------------------------------------------------------------------------

    function Cancel_Sell_Order() internal { // internal

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 41; // 41 = Cancel_Sell_Order

        int flag = 0;

        for (int i = 1; i <= No_Sell_Orders+1; i++){

            if (Sells[i].Price == Price && msg.sender == Sells[i].Address  &&  flag == 0) {

                if (Own_Amount[msg.sender] < 0){

                    Own_Amount_Sell_Order[msg.sender] -= Sells[i].Amount * Sells[i].Security * sU / (Sells[i].Price * Sells[i].Amount);
                    Own_Funds[msg.sender]             += Sells[i].Security;

                }

                Own_Amount[msg.sender] += Sells[i].Amount;
                No_Sell_Orders -= 1;
                flag = 1;
            }

            if (flag == 1) {
                Sells[i].Price      = Sells[i+1].Price;
                Sells[i].Amount     = Sells[i+1].Amount;
                Sells[i].Address    = Sells[i+1].Address;
                Sells[i].Security   = Sells[i+1].Security;
            }
        }

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 90; // 90 = exit
        Own_FeedBack[msg.sender] += flag;           // 91 canceled / 90 not canceled
    }

    function Cancel_Buy_Order() internal { // internal

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 42; // 42 = Cancel_Buy_Order

        int flag = 0;

        for (int i = 1; i <= No_Buy_Orders+1; i++){

            if (Buys[i].Price == Price && msg.sender == Buys[i].Address  &&  flag == 0) {

                Own_Funds[msg.sender] += Buys[i].Price * Buys[i].Amount / sU;
                No_Buy_Orders -= 1;
                flag = 1;

            }

            if (flag == 1) {
                Buys[i].Price      = Buys[i+1].Price;
                Buys[i].Amount     = Buys[i+1].Amount;
                Buys[i].Address    = Buys[i+1].Address;
            }

        }

        Own_FeedBack[msg.sender] = Own_FeedBack[msg.sender] * 100 + 90; // 90 = exit
        Own_FeedBack[msg.sender] += flag;           // 91 canceled / 90 not canceled
    }

// ------------------------------------------------------------------------------
// Reference Price
// ------------------------------------------------------------------------------
    function Set_Ref_Price(int List_Price) internal {
        if (List_Price > 0 && Ref_Price_Block_No + 100 > int(block.number)){
            Ref_Price = (Ref_Price * 99 + List_Price) / 100;
            Ref_Price_Block_No = int(block.number);
        }
    }
// ------------------------------------------------------------------------------
// universal routines
// ------------------------------------------------------------------------------
// https://github.com/ethereum/wiki/wiki/Solidity-Features

  function max(int a, int b) internal returns (int) { // internal
    if (a > b) return a;
    else return b;
  }
  function min(int a, int b) internal returns (int) { // internal
    if (a < b) return a;
    else return b;
  }

}
