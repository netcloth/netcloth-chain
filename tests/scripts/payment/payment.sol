pragma solidity 0.6.0;

import "./math/SafeMath.sol";

contract Payment {
    using SafeMath for uint256;

    event Transfer(address _from, address _to, uint256 _value);

    // 佣金地址
    address payable public _beneficiary;
    // 交易额费率，默认万分之一 (1 / 10**4)
    uint256 public _feeRateE4 = 1;

    // 构造函数，传入佣金地址
    constructor (address payable beneficiary) public {
        require(beneficiary != address(0), "beneficiary is the zero address");
        _beneficiary = beneficiary;
    }

    // 设置费率
    function setFeeRate(uint256 _feeRate) public {
        require(_beneficiary == msg.sender, "only beneficiary can set feeRate");
        _feeRateE4 = _feeRate;
    }

    // 结算
    function payAndStore(address payable _payee, uint256 _orderId, bytes32 _hash, uint256 expireTimeSec,  bytes32 R, bytes32 S, uint8 V) public payable {
        uint256 value = msg.value;
        
        bytes memory d = abi.encodePacked(_payee, _orderId, _hash, expireTimeSec);
        bytes32 hash = sha256(d);
        address expected_addr = ecrecover(hash, V, R, S);
        require (expected_addr == _payee, "signature verify falied!");
        
        uint256 commission = value.mul(_feeRateE4).div(uint256(10000));
        uint256 amount = value.sub(commission);
        _payee.transfer(amount);
        emit Transfer(msg.sender, _payee, amount);
    }
    
    // 佣金提现
    function withdraw(uint256 amount) public {
        require (msg.sender != _beneficiary);
        _beneficiary.transfer(amount);
    }
}
