package test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func (suite *TestSuite) TestAPI_GetTxs_Successfully() {
	name := "TestAPI_GetTxs_Successfully"
	suite.T().Log(name)

	// mock db
	db.Exec(`INSERT INTO public.pools (id, address, start_block, token0, token0_decimals, token1, token1_decimals, created_at) VALUES(1, '0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640', 19639191, '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48', 6, '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2', 18, 1712920902);`)
	db.Exec(`INSERT INTO public.txs (id, pool_id, tx_hash, block_number, block_time, gas, gas_price, receipt, is_finalized, created_at) VALUES(1, 1, '0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f', 19639195, 1712920871, 351800, 16089958436, '{"type":"0x2","root":"0x","status":"0x1","cumulativeGasUsed":"0x3792cf","logsBloom":"0x00200000010000400000000080000000000000000000000000000000040000000000200000000000000088000000000002000000880020000000000001200000000000080000000808004008000000200000020000000800000000000020102000000000100002000000000000000000000000000000000000082010040880000000000000000000000000000000000000000000010000080000004000000000020000004000200000000000000000020000000000000000002000000008000000002002000000000080000000000008000000000000009000000000000000080010200000000000000010000000008000001000000000002010000000004000","logs":[{"address":"0xc5fb36dd2fb59d3b98deff88425a3f425ee469ed","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x000000000000000000000000418503a66f9aed73afddea68827dd2ee8997495c","0x00000000000000000000000067cea36eeb36ace126a3ca6e21405258130cf33c"],"data":"0x00000000000000000000000000000000000000000000000000001476b081e800","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x56","removed":false},{"address":"0xc5fb36dd2fb59d3b98deff88425a3f425ee469ed","topics":["0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925","0x000000000000000000000000418503a66f9aed73afddea68827dd2ee8997495c","0x000000000000000000000000000000000022d473030f116ddee9f6b43ac78ba3"],"data":"0xffffffffffffffffffffffffffffffffffffffffffffffffffff509c8bab52f2","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x57","removed":false},{"address":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x00000000000000000000000067cea36eeb36ace126a3ca6e21405258130cf33c","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad"],"data":"0x000000000000000000000000000000000000000000000000000000002e75f7f7","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x58","removed":false},{"address":"0x67cea36eeb36ace126a3ca6e21405258130cf33c","topics":["0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1"],"data":"0x000000000000000000000000000000000000000000000000000000ca9a2e6b2400000000000000000000000000000000000000000000000000590c59c6d2824d","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x59","removed":false},{"address":"0x67cea36eeb36ace126a3ca6e21405258130cf33c","topics":["0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad"],"data":"0x000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001476b081e800000000000000000000000000000000000000000000000000000000002e75f7f70000000000000000000000000000000000000000000000000000000000000000","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x5a","removed":false},{"address":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x00000000000000000000000088e6a0c2ddd26feeb64f039a2c41296fcb3f5640","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad"],"data":"0x000000000000000000000000000000000000000000000000030ec4e06822b66c","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x5b","removed":false},{"address":"0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad","0x00000000000000000000000088e6a0c2ddd26feeb64f039a2c41296fcb3f5640"],"data":"0x000000000000000000000000000000000000000000000000000000002e75f7f7","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x5c","removed":false},{"address":"0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640","topics":["0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad"],"data":"0x000000000000000000000000000000000000000000000000000000002e75f7f7fffffffffffffffffffffffffffffffffffffffffffffffffcf13b1f97dd499400000000000000000000000000000000000041b0be4c4d3dbe3ae2220b77e8fe0000000000000000000000000000000000000000000000026a832d11fa3e5a9e000000000000000000000000000000000000000000000000000000000002f834","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x5d","removed":false},{"address":"0x32c6f1c1731ff8f98ee2ede8954f696446307846","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x000000000000000000000000bd35130bd84a3f016d1b6af257e0e0544887a42a","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad"],"data":"0x000000000000000000000000000000000000000000002001644e714413e2e300","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x5e","removed":false},{"address":"0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad","0x000000000000000000000000bd35130bd84a3f016d1b6af257e0e0544887a42a"],"data":"0x000000000000000000000000000000000000000000000000030ec4e06822b66c","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x5f","removed":false},{"address":"0xbd35130bd84a3f016d1b6af257e0e0544887a42a","topics":["0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad"],"data":"0xffffffffffffffffffffffffffffffffffffffffffffdffe9bb18ebbec1d1d00000000000000000000000000000000000000000000000000030ec4e06822b66c0000000000000000000000000000000000000000004efee2c0608896778fe06e0000000000000000000000000000000000000000000005b74fabc7170da930b1fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffdf2e6","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x60","removed":false},{"address":"0x32c6f1c1731ff8f98ee2ede8954f696446307846","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad","0x00000000000000000000000037a8f295612602f2774d331e562be9e61b83a327"],"data":"0x0000000000000000000000000000000000000000000000147bc550ec54879b72","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x61","removed":false},{"address":"0x32c6f1c1731ff8f98ee2ede8954f696446307846","topics":["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef","0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad","0x000000000000000000000000418503a66f9aed73afddea68827dd2ee8997495c"],"data":"0x000000000000000000000000000000000000000000001fece8892057bf5b478e","blockNumber":"0x12bab9b","transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","transactionIndex":"0x34","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","logIndex":"0x62","removed":false}],"transactionHash":"0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f","contractAddress":"0x0000000000000000000000000000000000000000","gasUsed":"0x55e38","effectiveGasPrice":"0x3bf094824","blockHash":"0xd3ffa8172bebc0c064514da3907f4737cfbd6b9c419f85a0c68ee75bd8dbea49","blockNumber":"0x12bab9b","transactionIndex":"0x34"}'::json, false, 1712920927);`)
	db.Exec(`INSERT INTO public.swap_events (id, pool_id, tx_id, amount0, amount1, price, created_at) VALUES(1, 1, 1, '779483127', '-220329899886556780', '3537.800033', 1712920927);`)
	db.Exec(`INSERT INTO public.ethusdt_klines (id, open_time, close_time, open_price, high_price, low_price, close_price, ohlc4, created_at) VALUES(1, 1712920860000, 1712920919999, '3537.00000000', '3538.80000000', '3537.00000000', '3538.79000000', 3537.8975000000, 1712920928);`)

	// call API
	method := "GET"
	endpoint := fmt.Sprintf("/api/v1/txs?poolAddress=%s&pageSize=1", "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640")

	/// call
	req, _ := http.NewRequest(method, endpoint, nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert response
	suite.EqualValues(http.StatusOK, w.Code)
	actualResBody := map[string]interface{}{}
	json.Unmarshal(responseData, &actualResBody)
	expected := map[string]interface{}{
		"code": float64(0),
		"data": map[string]interface{}{
			"pagination": map[string]interface{}{
				"total": float64(1),
			},
			"txs": []interface{}{
				map[string]interface{}{
					"blockNumber": float64(19639195),
					"blockTime":   float64(1712920871),
					"createdAt":   float64(1712920927),
					"gas":         float64(351800),
					"gasPrice":    "16089958436",
					"id":          float64(1),
					"isFinalized": false,
					"poolId":      float64(1),
					"swapEvents": []interface{}{
						map[string]interface{}{
							"amount0":   "779483127",
							"amount1":   "-220329899886556780",
							"createdAt": float64(1712920927),
							"id":        float64(1),
							"price":     "3537.800033",
						},
					},
					"txHash": "0xc5f30ba6a4dcfa616525ebf3b1a64396b0d1f27341e9c0998e7347bad9d3da8f",
				},
			},
		},
		"message":   "successfully",
		"requestId": "",
	}
	suite.EqualValues(expected, actualResBody)
}

func (suite *TestSuite) TestAPI_GetTxs_Failed_1() {
	name := "TestAPI_GetTxs_Failed_1: invalid pool address"
	suite.T().Log(name)

	// call API
	method := "GET"
	endpoint := fmt.Sprintf("/api/v1/txs?poolAddress=%s&pageSize=1", "0x_invalid")

	/// call
	req, _ := http.NewRequest(method, endpoint, nil)
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	actualResBody := map[string]interface{}{}
	json.Unmarshal(responseData, &actualResBody)
	expected := map[string]interface{}{
		"code": float64(4000),
		"details": []interface{}{
			map[string]interface{}{
				"fieldViolations": []interface{}{
					map[string]interface{}{"description": "invalid", "field": "poolAddress"},
				},
			},
		},
		"message":   "bad request",
		"requestId": "",
	}
	suite.EqualValues(expected, actualResBody)
}
