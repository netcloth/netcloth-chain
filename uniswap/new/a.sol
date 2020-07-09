pragma solidity 0.6.0;

contract Account{
    uint accId;
}

contract Test{
    address public a;

    constructor() public {
        Account b = new Account();
        a = address(b);
    }

    function  newAccount() public {
        Account b = new Account();
        a = address(b);
    }

}