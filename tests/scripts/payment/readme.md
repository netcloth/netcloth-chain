# payment合约

支付和存证合约，合约创建者可设置交易佣金收益地址和费率。

## 接口

### constructor
构造函数，传入佣金地址

### setFeeRate

设置交易费率，对单笔交易按交易额抽成，默认为万分之1。

```javascript
function setFeeRate(uint256 _feeRate) public
```

### payAndStore

支付和存证接口，用户调用合约时，传入收款人地址和订单信息及签名。

调用成功，收款人直接到账，合约不托管资产。

```javascript
// _payee 收款人地址
// _orderId 订单id
// _hash  商品摘要信息
// expireTimeSec 失效时间，单位为秒
// R,S,V 为收款人对上述字段的签名，合约内会验证，防篡改
function payAndStore(address payable _payee, uint256 _orderId, bytes32 _hash, uint256 expireTimeSec,  bytes32 R, bytes32 S, uint8 V) public payable
```

调用成功，触发事件

```javascript
event Transfer(address _from, address _to, uint256 _value);
```

### withdraw

取回佣金。