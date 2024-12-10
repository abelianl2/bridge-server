package server

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_SaveTxAndMemo(t *testing.T) {

	//powerOfTen := new(big.Float).SetFloat64(1e7)
	//
	//x := new(big.Float).SetFloat64(2000)
	//
	//// 计算浮点数乘以10的7次方
	//r := new(big.Float).Mul(x, powerOfTen)
	//
	str := base64.StdEncoding.EncodeToString([]byte("20000000000"))
	//
	//t.Log(str)

	//str := `ABABrkOx//97InByb3RvY29sIjoiTWFibGUiLCJhY3Rpb24iOiJkZXBvc2l0IiwiZnJvbSI6ImFiZTM2ZjUwM2UxNGY5ZmUxMzk1MGUwMDlkODlkZTI2OTAzMWFhYjA1NDIyMzg1OGNjNDI0MTIyNGI5NWM5ZmQwYmVkMzgxZDQ0NWNhMTA3N2I2OWY0YmQxMmZhYTIyNDg3OTdmNmVkYWVlN2Q0Nzc3ZmYxYTYzNjZmM2E0NmQxOThkOCIsInJlY2VpcHQiOiIweGRhYzE3Zjk1OGQyZWU1MjNhMjIwNjIwNjk5NDU5N2MxM2Q4MzFlYzciLCJ0byI6ImFiZTMzOGNlMGNlMTc4ZmIwYWNhNDJiNGU0MDBjZGYzOTVjOTJjYmY5YzVjOWFiZDY3OGFhNTE2ODM1ZjY5N2JkNmQyODBiMjg1ODE1OTI0Zjg2MjM1MmM1NDYzNDIxYzlmOGQyNDdmNjVkYzExMmFhMDRjMjVkZTkyNWJkMWQxYTMzNCIsInZhbHVlIjoiMTEwMDAwMDAiLCJsb2NrdXBQZXJpb2QiOjE4MCwicmV3YXJkUmF0aW8iOjF9`
	bs, err := base64.StdEncoding.DecodeString(str)

	assert.NoError(t, err)

	t.Log(string(bs))
}
