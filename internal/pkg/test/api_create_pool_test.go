package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/golang/mock/gomock"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	"github.com/vuquang23/poseidon/internal/pkg/server/middleware"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	"github.com/vuquang23/poseidon/pkg/asynq"
)

func (suite *TestSuite) TestAPI_CreatePool_Successfully() {
	name := "TestAPI_CreatePool_Successfully: create pool successfully"
	suite.T().Log(name)

	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockAsynqClient := asynq.NewMockIAsynqClient(ctrl)
	poolRepo.SetAsynqClient(mockAsynqClient)

	mockAsynqClient.EXPECT().EnqueueTask(
		gomock.Any(),
		string(valueobject.TaskTypeHandlePoolCreated),
		"", "", gomock.Any(), -1,
	).Return(nil).Times(1)

	// call API
	method := "POST"
	endpoint := "/api/v1/pools"

	/// init req body
	body := map[string]interface{}{
		"address":        "0x88E6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		"token0":         "0x88E6a0c2ddd26feeb64f039a2c41296fcb3f5641",
		"token0Decimals": 18,
		"token1":         "0x88E6a0c2ddd26feeb64f039a2c41296fcb3f5642",
		"token1Decimals": 6,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(bodyBytes))
	req.Header.Add(middleware.HeaderXAPIKey, XAPIKey)

	/// call
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert database
	var pool entity.Pool
	db.First(&pool)
	suite.EqualValues("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640", pool.Address)
	suite.EqualValues("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5641", pool.Token0)
	suite.EqualValues(18, pool.Token0Decimals)
	suite.EqualValues("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5642", pool.Token1)
	suite.EqualValues(6, pool.Token1Decimals)

	// assert API response
	suite.EqualValues(http.StatusOK, w.Code)
	actualResBody := map[string]interface{}{}
	json.Unmarshal(responseData, &actualResBody)
	suite.Equal(float64(0), actualResBody["code"])
	suite.Equal("successfully", actualResBody["message"])
	suite.Positive(actualResBody["data"].(map[string]interface{})["poolId"])
}

func (suite *TestSuite) TestAPI_CreatePool_Failed() {
	name := "TestAPI_CreatePool_Failed: invalid pool address"
	suite.T().Log(name)

	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockAsynqClient := asynq.NewMockIAsynqClient(ctrl)
	poolRepo.SetAsynqClient(mockAsynqClient)

	mockAsynqClient.EXPECT().EnqueueTask(
		gomock.Any(),
		string(valueobject.TaskTypeHandlePoolCreated),
		"", "", gomock.Any(), -1,
	).Return(nil).Times(0)

	// call API
	method := "POST"
	endpoint := "/api/v1/pools"

	/// init req body
	body := map[string]interface{}{
		"address":        "0x_invalid",
		"token0":         "0x88E6a0c2ddd26feeb64f039a2c41296fcb3f5641",
		"token0Decimals": 18,
		"token1":         "0x88E6a0c2ddd26feeb64f039a2c41296fcb3f5642",
		"token1Decimals": 6,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, endpoint, bytes.NewReader(bodyBytes))
	req.Header.Add(middleware.HeaderXAPIKey, XAPIKey)

	/// call
	w := httptest.NewRecorder()
	ginEngine.ServeHTTP(w, req)
	responseData, _ := io.ReadAll(w.Body)

	// assert API response
	suite.EqualValues(http.StatusBadRequest, w.Code)
	actualResBody := map[string]interface{}{}
	json.Unmarshal(responseData, &actualResBody)
	suite.Equal(float64(4000), actualResBody["code"])
	suite.Equal("bad request", actualResBody["message"])
}
