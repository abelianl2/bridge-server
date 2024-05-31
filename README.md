# bridge-server

The backend service code of the bridge portal

### usage


1. 发送deposit交易请求

`````` 
curl --location --request POST '190.92.213.101:9002/api/bridge/submitWithMemo' \
--header 'Content-Type: application/json' \
--data-raw '{
        "from_network": "abelian-test", // from network
        "from_address":"abe36f503e14f9fe13950e009d89de269031aab054223858cc4241224b95c9fd0bed381d445ca1077b69f4bd12faa2248797f6edaee7d4777ff1a6366f3a46d198d8", //L1 from address
        "to_network": "mable-test", //to network
        "to_address": "0xdac17f958d2ee523a2206206994597c13d831ec7", //L2 recipient address
        "amount":100000
    }'


response:

{
    "code": 0,
    "data": "http://190.92.213.101:9002/api/bridge/deposit/73659d56-3b0e-42c1-a1d5-1b5c0e466c6c",
    "message": "ok"
}

``````

2. 钱包扫描

 钱包通过扫描 bridge 提供的二维码, 得到 一个连接数据，如下：

``````
 http://190.92.213.101:9002/api/bridge/deposit/73659d56-3b0e-42c1-a1d5-1b5c0e466c6c

``````

3. 请求deposit交易所需要数据

`````` 
curl --location --request GET 'http://190.92.213.101:9002/api/bridge/deposit/21020787-4a2b-40c1-9bd3-98a034f83102'

response:

{
    "code": 0,
    "data": {
        "amountOfGasFee": "0.09",
        "callback": "http://190.92.213.101:9002/api/bridge/notify/21020787-4a2b-40c1-9bd3-98a034f83102",
        "recipient": "abe338ce0ce178fb0aca42b4e400cdf395c92cbf9c5c9abd678aa516835f697bd6d280b285815924f862352c5463421c9f8d247f65dc112aa04c25de925bd1d1a334",
        "senderMd5": "86314525cbb769d9d7b40d9e58f67a13",
        "memo": "ABABiWQG//97InByb3RvY29sIjoiTWFibGUiLCJhY3Rpb24iOiJkZXBvc2l0IiwiZnJvbSI6ImFiZTM2ZjUwM2UxNGY5ZmUxMzk1MGUwMDlkODlkZTI2OTAzMWFhYjA1NDIyMzg1OGNjNDI0MTIyNGI5NWM5ZmQwYmVkMzgxZDQ0NWNhMTA3N2I2OWY0YmQxMmZhYTIyNDg3OTdmNmVkYWVlN2Q0Nzc3ZmYxYTYzNjZmM2E0NmQxOThkOCIsInJlY2VpcHQiOiIweGRhYzE3Zjk1OGQyZWU1MjNhMjIwNjIwNjk5NDU5N2MxM2Q4MzFlYzciLCJ0byI6ImFiZTMzOGNlMGNlMTc4ZmIwYWNhNDJiNGU0MDBjZGYzOTVjOTJjYmY5YzVjOWFiZDY3OGFhNTE2ODM1ZjY5N2JkNmQyODBiMjg1ODE1OTI0Zjg2MjM1MmM1NDYzNDIxYzlmOGQyNDdmNjVkYzExMmFhMDRjMjVkZTkyNWJkMWQxYTMzNCIsInZhbHVlIjoiMTAwMDAwIn0=",
        "amountOfAbel": "100000"
    },
    "message": "ok"
}

``````

4. 向链上提交交易

 钱包首先验证 senderMd5 是否合法，验证通过后，向链上发送交易并得到 TxId

5. 通知调用方
 
 根据 步骤3 中得到 callback url ,通知调用方 交易状态

``````  

curl --location --request POST '190.92.213.101:9002/api/bridge/notify/21020787-4a2b-40c1-9bd3-98a034f83102' \
--header 'Content-Type: application/json' \
--data-raw '{
        "hash": "0x1233333"
    }'
    
  
resonse:

{
    "code": 0,
    "data": null,
    "message": "ok"
}     


``````







