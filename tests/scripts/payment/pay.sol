pragma solidity 0.6.0;

import "./SafeMath.sol";

contract Pay {
    
    using SafeMath for uint256;
    
    event Transfer(address _from, address _to, uint256 _value, uint256 _actual_value);
    
    uint256 public E4 = uint256(10000);
    uint256 public feeRateE4;
    address payable public owner;
    
    constructor(uint256 _feeRateE4) public {
        assert(feeRateValid(_feeRateE4));
        owner = msg.sender;
        feeRateE4 = _feeRateE4;
    }
    
    function setFeeRate(uint256 _feeRateE4) public {
        assert(msg.sender == owner);
        assert(feeRateValid(_feeRateE4));
        
        feeRateE4 = _feeRateE4;
    }
    
    function feeRateValid(uint256 _feeRateE4) public view returns(bool) {
        return _feeRateE4 <= E4;
    }
    
    function calcCommission(uint256 value) public view returns(uint256) {
        return value.mul(feeRateE4).div(E4);
    }
    
    function doTransfer(address payable to) public payable {
        uint256 commission = calcCommission(msg.value);
        uint256 valueToSendout = msg.value - commission;
        to.transfer(valueToSendout);
        emit Transfer(msg.sender, to, msg.value, valueToSendout);
    }
    
	function withdraw(uint256 amount) public {
		assert (msg.sender == owner);
		owner.transfer(amount);
	}
}

