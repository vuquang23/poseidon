package test

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	"github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/eth"
)

func (suite *TestSuite) TestTask_ScanTxs_Successfully() {
	name := "TestTask_ScanTxs_Successfully"
	suite.T().Log(name)

	var (
		poolAddress        = common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640")
		fromBlock   uint64 = 12379656
		toBlock     uint64 = 12379659
	)

	// mock DB
	db.Exec(`INSERT INTO public.pools (id, address, start_block, token0, token0_decimals, token1, token1_decimals, created_at) VALUES(1, '0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640', 12379656, '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48', 6, '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2', 18, 1712856369);`)
	db.Exec(`INSERT INTO public.block_cursors (id, pool_id, "type", block_number, extra, created_at, updated_at) VALUES(1, 1, 'scanner', 12379656, NULL, 1712856369, 1712856412);`)

	// mock call
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockEthClient := eth.NewMockIClient(ctrl)
	mockEthClient.EXPECT().
		GetLatestBlockHeader(gomock.Any()).
		Return(&types.Header{Number: big.NewInt(int64(toBlock))}, nil).
		Times(1)

	mockEthClient.EXPECT().
		GetLogs(gomock.Any(), fromBlock, toBlock, []common.Address{poolAddress}).
		Return([]types.Log{
			{
				Address: poolAddress,
				Topics: []common.Hash{
					common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"),
					common.HexToHash("0x000000000000000000000000e592427a0aece92de3edee1f18e0157c05861564"),
					common.HexToHash("0x0000000000000000000000003cec6746ebd7658f58e5d786e0999118fea2905c"),
				},
				Data:        common.FromHex("0x000000000000000000000000000000000000000000000000000000003aca2d71fffffffffffffffffffffffffffffffffffffffffffffffffc0693922f3e7040000000000000000000000000000000000000427f3dbad735381bf74afafab87d00000000000000000000000000000000000000000000000000189fc17963f9e1000000000000000000000000000000000000000000000000000000000002f928"),
				BlockNumber: 12379656,
				TxHash:      common.HexToHash("0x6594e6beb27a2bd0ea0d23599ffc4343aeab961438c5d1c5a38d53ae0431daf2"),
				TxIndex:     0,
				BlockHash:   common.HexToHash("0x34adfd9af10e6cd70a4afad80cf0c571bbe941e1e6894d80420d2517f5c7c4b1"),
				Index:       1,
				Removed:     false,
			},
		}, nil).
		Times(1)

	mockEthClient.EXPECT().
		HeaderByHash(gomock.Any(), common.HexToHash("0x34adfd9af10e6cd70a4afad80cf0c571bbe941e1e6894d80420d2517f5c7c4b1")).
		Return(&types.Header{
			ParentHash:       common.HexToHash("0xa8ea0cf850e775dbb06ca03cd5163458a18bddb2405e3d5d9d524e870a19c716"),
			UncleHash:        common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
			Coinbase:         common.HexToAddress("0x02ad7c55a19e976ec105172a75a9d84dc9cf23c6"),
			Root:             common.HexToHash("0x30b54f4d718a328f9d409280ca4f354b4e27888a45de8c2fe44d59b7edcaf568"),
			TxHash:           common.HexToHash("0x12fec646d6730b2c100925c6113a3d3702e62b7e69eee342cebac977fe532413"),
			ReceiptHash:      common.HexToHash("0x871d0ec3bcb97f1893ffcbca2d6bd93a101cdc067a785823abcdee48f1b168dd"),
			Bloom:            types.BytesToBloom(common.FromHex("0xf8ee7833f90cc79750923d798b571ef81de6dce932d56028b6dd754e673797cb81c5c553146a383b551a7d02924a3d55422bd2248b03ffffa633dbe03370800c718960c83792fb4be9d6062f1fc1bae93f9871235a63234538067fc7dc4e7a26fbc8a38f9e9ecdcd191e5fd8cdd2fff32f6c2c71e9fdd4f06646e1df368e72e028b793ab9bac7efc8ecd45dcb1adeb4cc597d879e7f8247dc5ae2a4f019b88dd77ffb82f2e1665d21eee0495a9eb234567a1a58518678abf11b32a16198c0f347d548536454a00f5e82b1d30a742d62cdb38579ed492fbd222117293642d6a95f8bd2d04c09ddcac27a5d5e833db18bc06f86a506e573dded82085bcd063dd85")),
			Difficulty:       big.NewInt(7618240087310300),
			Number:           big.NewInt(12379656),
			GasLimit:         14977980,
			GasUsed:          14973863,
			Time:             1620289565,
			Extra:            common.FromHex("0xd883010a02846765746888676f312e31332e34856c696e7578"),
			MixDigest:        common.HexToHash("0x7f8e4e207024f1dae9406eb0a6f8c07152cfde4b8b618b67387bd31b5d751f84"),
			Nonce:            types.EncodeNonce(0x98bc747da448fb79),
			BaseFee:          nil,
			WithdrawalsHash:  nil,
			BlobGasUsed:      nil,
			ExcessBlobGas:    nil,
			ParentBeaconRoot: nil,
		}, nil).
		Times(1)

	mockEthClient.EXPECT().
		GetTxReceipt(gomock.Any(), common.HexToHash("0x6594e6beb27a2bd0ea0d23599ffc4343aeab961438c5d1c5a38d53ae0431daf2")).
		Return(&types.Receipt{
			TxHash:            common.HexToHash("0x6594e6beb27a2bd0ea0d23599ffc4343aeab961438c5d1c5a38d53ae0431daf2"),
			BlockHash:         common.HexToHash("0x34adfd9af10e6cd70a4afad80cf0c571bbe941e1e6894d80420d2517f5c7c4b1"),
			GasUsed:           100000,
			EffectiveGasPrice: big.NewInt(11119999999),
		}, nil)

	mockAsynqClient := asynq.NewMockIAsynqClient(ctrl)
	mockAsynqClient.EXPECT().EnqueueTask(
		gomock.Any(), valueobject.TaskTypeGetETHUSDTKline, "", "", valueobject.TaskGetETHUSDTKlinePayload{Time: 1620289560}, -1,
	).Return(nil).Times(1)

	taskSvc.SetEthClient(mockEthClient)
	taskSvc.SetAsynqClient(mockAsynqClient)

	// handle
	payload := valueobject.TaskScanTxsPayload{
		PoolID:         1,
		PoolAddress:    "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		Token0Decimals: 6,
		Token1Decimals: 18,
	}

	err := taskSvc.ScanTxs(context.Background(), payload)
	suite.Nil(err)

	// assert DB
	var cursor entity.BlockCursor
	db.First(&cursor)
	suite.EqualValues(12379660, cursor.BlockNumber)

	var tx entity.Tx
	db.First(&tx)
	suite.EqualValues("0x6594e6beb27a2bd0ea0d23599ffc4343aeab961438c5d1c5a38d53ae0431daf2", tx.TxHash)
	suite.EqualValues(12379656, tx.BlockNumber)
	suite.EqualValues(1620289565, tx.BlockTime)
	suite.EqualValues(100000, tx.Gas)
	suite.True(tx.GasPrice.Equal(decimal.NewFromInt(11119999999)))
	suite.False(tx.IsFinalized)

	var event entity.SwapEvent
	db.First(&event)
	suite.EqualValues(tx.ID, event.TxID)
	suite.EqualValues("986328433", event.Amount0)
	suite.EqualValues("-286379270224318400", event.Amount1)
	suite.EqualValues("3444.133482", event.Price)
}

func (suite *TestSuite) TestTask_ScanTxs_Failed() {
	name := "TestTask_ScanTxs_Failed: can not get tx receipt"
	suite.T().Log(name)

	var (
		poolAddress        = common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640")
		fromBlock   uint64 = 12379656
		toBlock     uint64 = 12379659

		errGetReceiptFailed = errors.New("get receipt failed")
	)

	// mock DB
	db.Exec(`INSERT INTO public.pools (id, address, start_block, token0, token0_decimals, token1, token1_decimals, created_at) VALUES(1, '0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640', 12379656, '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48', 6, '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2', 18, 1712856369);`)
	db.Exec(`INSERT INTO public.block_cursors (id, pool_id, "type", block_number, extra, created_at, updated_at) VALUES(1, 1, 'scanner', 12379656, NULL, 1712856369, 1712856412);`)

	// mock eth call
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockEthClient := eth.NewMockIClient(ctrl)
	mockEthClient.EXPECT().
		GetLatestBlockHeader(gomock.Any()).
		Return(&types.Header{Number: big.NewInt(int64(toBlock))}, nil).
		Times(1)

	mockEthClient.EXPECT().
		GetLogs(gomock.Any(), fromBlock, toBlock, []common.Address{poolAddress}).
		Return([]types.Log{
			{
				Address: poolAddress,
				Topics: []common.Hash{
					common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"),
					common.HexToHash("0x000000000000000000000000e592427a0aece92de3edee1f18e0157c05861564"),
					common.HexToHash("0x0000000000000000000000003cec6746ebd7658f58e5d786e0999118fea2905c"),
				},
				Data:        common.FromHex("0x000000000000000000000000000000000000000000000000000000003aca2d71fffffffffffffffffffffffffffffffffffffffffffffffffc0693922f3e7040000000000000000000000000000000000000427f3dbad735381bf74afafab87d00000000000000000000000000000000000000000000000000189fc17963f9e1000000000000000000000000000000000000000000000000000000000002f928"),
				BlockNumber: 12379656,
				TxHash:      common.HexToHash("0x6594e6beb27a2bd0ea0d23599ffc4343aeab961438c5d1c5a38d53ae0431daf2"),
				TxIndex:     0,
				BlockHash:   common.HexToHash("0x34adfd9af10e6cd70a4afad80cf0c571bbe941e1e6894d80420d2517f5c7c4b1"),
				Index:       1,
				Removed:     false,
			},
		}, nil).
		Times(1)

	mockEthClient.EXPECT().
		HeaderByHash(gomock.Any(), common.HexToHash("0x34adfd9af10e6cd70a4afad80cf0c571bbe941e1e6894d80420d2517f5c7c4b1")).
		Return(&types.Header{
			ParentHash:       common.HexToHash("0xa8ea0cf850e775dbb06ca03cd5163458a18bddb2405e3d5d9d524e870a19c716"),
			UncleHash:        common.HexToHash("0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"),
			Coinbase:         common.HexToAddress("0x02ad7c55a19e976ec105172a75a9d84dc9cf23c6"),
			Root:             common.HexToHash("0x30b54f4d718a328f9d409280ca4f354b4e27888a45de8c2fe44d59b7edcaf568"),
			TxHash:           common.HexToHash("0x12fec646d6730b2c100925c6113a3d3702e62b7e69eee342cebac977fe532413"),
			ReceiptHash:      common.HexToHash("0x871d0ec3bcb97f1893ffcbca2d6bd93a101cdc067a785823abcdee48f1b168dd"),
			Bloom:            types.BytesToBloom(common.FromHex("0xf8ee7833f90cc79750923d798b571ef81de6dce932d56028b6dd754e673797cb81c5c553146a383b551a7d02924a3d55422bd2248b03ffffa633dbe03370800c718960c83792fb4be9d6062f1fc1bae93f9871235a63234538067fc7dc4e7a26fbc8a38f9e9ecdcd191e5fd8cdd2fff32f6c2c71e9fdd4f06646e1df368e72e028b793ab9bac7efc8ecd45dcb1adeb4cc597d879e7f8247dc5ae2a4f019b88dd77ffb82f2e1665d21eee0495a9eb234567a1a58518678abf11b32a16198c0f347d548536454a00f5e82b1d30a742d62cdb38579ed492fbd222117293642d6a95f8bd2d04c09ddcac27a5d5e833db18bc06f86a506e573dded82085bcd063dd85")),
			Difficulty:       big.NewInt(7618240087310300),
			Number:           big.NewInt(12379656),
			GasLimit:         14977980,
			GasUsed:          14973863,
			Time:             1620289565,
			Extra:            common.FromHex("0xd883010a02846765746888676f312e31332e34856c696e7578"),
			MixDigest:        common.HexToHash("0x7f8e4e207024f1dae9406eb0a6f8c07152cfde4b8b618b67387bd31b5d751f84"),
			Nonce:            types.EncodeNonce(0x98bc747da448fb79),
			BaseFee:          nil,
			WithdrawalsHash:  nil,
			BlobGasUsed:      nil,
			ExcessBlobGas:    nil,
			ParentBeaconRoot: nil,
		}, nil).
		Times(1)

	mockEthClient.EXPECT().
		GetTxReceipt(gomock.Any(), common.HexToHash("0x6594e6beb27a2bd0ea0d23599ffc4343aeab961438c5d1c5a38d53ae0431daf2")).
		Return(nil, errGetReceiptFailed)

	mockAsynqClient := asynq.NewMockIAsynqClient(ctrl)
	mockAsynqClient.EXPECT().EnqueueTask(
		gomock.Any(), valueobject.TaskTypeGetETHUSDTKline, "", "", valueobject.TaskGetETHUSDTKlinePayload{Time: 1620289565}, -1,
	).Return(nil).Times(1)

	taskSvc.SetEthClient(mockEthClient)
	taskSvc.SetAsynqClient(mockAsynqClient)

	// handle
	payload := valueobject.TaskScanTxsPayload{
		PoolID:         1,
		PoolAddress:    "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		Token0Decimals: 6,
		Token1Decimals: 18,
	}

	err := taskSvc.ScanTxs(context.Background(), payload)
	suite.ErrorIs(err, errGetReceiptFailed)

	// assert DB
	var cursor entity.BlockCursor
	db.First(&cursor)
	suite.EqualValues(12379656, cursor.BlockNumber)

	var tx entity.Tx
	err = db.First(&tx).Error
	suite.ErrorIs(err, gorm.ErrRecordNotFound)

	var event entity.SwapEvent
	err = db.First(&event).Error
	suite.ErrorIs(err, gorm.ErrRecordNotFound)
}
