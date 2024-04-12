package test

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"

	"github.com/vuquang23/poseidon/internal/pkg/entity"
	tasksvc "github.com/vuquang23/poseidon/internal/pkg/service/task"
	"github.com/vuquang23/poseidon/internal/pkg/valueobject"
	"github.com/vuquang23/poseidon/pkg/asynq"
	"github.com/vuquang23/poseidon/pkg/eth"
)

func (suite *TestSuite) TestTask_FinalizeTxs_Successfully_1() {
	name := "TestTask_FinalizeTxs_Successfully_1: no reorg"
	suite.T().Log(name)

	var (
		poolAddress = common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640")
	)

	// mock DB
	db.Exec(`INSERT INTO public.pools (id, address, start_block, token0, token0_decimals, token1, token1_decimals, created_at) VALUES(1, '0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640', 19638503, '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48', 6, '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2', 18, 1712913866);`)
	db.Exec(`INSERT INTO public.block_cursors (id, pool_id, "type", block_number, extra, created_at, updated_at) VALUES(2, 1, 'finalizer', 19638547, '{"createdAtFinalizedBlock":19638538}'::json, 1712913868, 1712913962);
	INSERT INTO public.block_cursors (id, pool_id, "type", block_number, extra, created_at, updated_at) VALUES(1, 1, 'scanner', 19638622, NULL, 1712913868, 1712913963);`)

	db.Exec(`INSERT INTO public.txs (id, pool_id, tx_hash, block_number, block_time, gas, gas_price, receipt, is_finalized, created_at) VALUES(97, 1, '0x72e52be2d7f03b0842f4cc93e0fc798f340f71039bde322be821be7c8e8e677f', 19638615, 1712913887, 248711, 24859996161, '{}', false, 1712913893);`)
	db.Exec(`INSERT INTO public.swap_events (id, pool_id, tx_id, amount0, amount1, price, created_at) VALUES(96, 1, 97, '1160580000', '-329388970022326866', '3523.433101', 1712913893);`)

	// mock call
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockEthClient := eth.NewMockIClient(ctrl)
	mockEthClient.EXPECT().
		GetLatestBlockHeader(gomock.Any()).
		Return(&types.Header{Number: big.NewInt(int64(19638722)), Time: uint64(time.Now().Unix())}, nil).
		Times(1)

	mockEthClient.EXPECT().
		GetLogs(gomock.Any(), uint64(19638547), uint64(19638622-1), []common.Address{poolAddress}).
		Return([]types.Log{
			{
				Address: poolAddress,
				Topics: []common.Hash{
					common.HexToHash("0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"),
					common.HexToHash("0x0000000000000000000000003fc91a3afd70395cd496c647d5a6cc9d4b2b7fad"),
					common.HexToHash("0x000000000000000000000000e866771b1606356d0fc5c505049f005851771d9b"),
				},
				Data:        common.FromHex("0x00000000000000000000000000000000000000000000000000000000452d0ba0fffffffffffffffffffffffffffffffffffffffffffffffffb6dc67a18e32dae00000000000000000000000000000000000041d2fe22430496ac8e8ed669d1810000000000000000000000000000000000000000000000026bab927f43401f97000000000000000000000000000000000000000000000000000000000002f85d"),
				BlockNumber: 19638615,
				TxHash:      common.HexToHash("0x72e52be2d7f03b0842f4cc93e0fc798f340f71039bde322be821be7c8e8e677f"),
				TxIndex:     218,
				BlockHash:   common.HexToHash("0xc24b17bb74a96facace20f97872b59b6257f45f5133cb00879c360ba8011b159"),
				Index:       273,
				Removed:     false,
			},
		}, nil).
		Times(1)

	taskSvc.SetEthClient(mockEthClient)

	// finalize txs
	payload := valueobject.TaskFinalizeTxsPayload{
		PoolID:         1,
		PoolAddress:    "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		Token0Decimals: 6,
		Token1Decimals: 18,
	}

	err := taskSvc.FinalizeTxs(context.Background(), payload)
	suite.Nil(err)

	// check db
	var tx entity.Tx
	db.First(&tx)
	suite.True(tx.IsFinalized)

	var event entity.SwapEvent
	err = db.First(&event).Error
	suite.Nil(err)
}

func (suite *TestSuite) TestTask_FinalizeTxs_Successfully_2() {
	name := "TestTask_FinalizeTxs_Successfully_2: reorg occurs"
	suite.T().Log(name)

	// mock DB
	var (
		poolAddress        = common.HexToAddress("0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640")
		fromBlock   uint64 = 12379656
		toBlock     uint64 = 12379659
	)

	// mock DB
	db.Exec(`INSERT INTO public.pools (id, address, start_block, token0, token0_decimals, token1, token1_decimals, created_at) VALUES(1, '0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640', 12379656, '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48', 6, '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2', 18, 1712913866);`)
	db.Exec(`INSERT INTO public.block_cursors (id, pool_id, "type", block_number, extra, created_at, updated_at) VALUES(2, 1, 'finalizer', 12379656, '{"createdAtFinalizedBlock":12379581}'::json, 1712913868, 1712913962);
	INSERT INTO public.block_cursors (id, pool_id, "type", block_number, extra, created_at, updated_at) VALUES(1, 1, 'scanner', 12379660, NULL, 1712913868, 1712913963);`)

	db.Exec(`INSERT INTO public.txs (id, pool_id, tx_hash, block_number, block_time, gas, gas_price, receipt, is_finalized, created_at) VALUES(97, 1, '0x72e52be2d7f03b0842f4cc93e0fc798f340f71039bde322be821be7c8e8e677f', 12379657, 1712913887, 248711, 24859996161, '{}', false, 1712913893);`)
	db.Exec(`INSERT INTO public.swap_events (id, pool_id, tx_id, amount0, amount1, price, created_at) VALUES(96, 1, 97, '1160580000', '-329388970022326866', '3523.433101', 1712913893);`)

	// mock call
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockEthClient := eth.NewMockIClient(ctrl)
	mockEthClient.EXPECT().
		GetLatestBlockHeader(gomock.Any()).
		Return(&types.Header{Number: big.NewInt(int64(fromBlock + 100)), Time: uint64(time.Now().Unix())}, nil).
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
	payload := valueobject.TaskFinalizeTxsPayload{
		PoolID:         1,
		PoolAddress:    "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		Token0Decimals: 6,
		Token1Decimals: 18,
	}

	err := taskSvc.FinalizeTxs(context.Background(), payload)
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
	suite.True(tx.IsFinalized)

	var event entity.SwapEvent
	db.First(&event)
	suite.EqualValues(tx.ID, event.TxID)
	suite.EqualValues("986328433", event.Amount0)
	suite.EqualValues("-286379270224318400", event.Amount1)
	suite.EqualValues("3444.133482", event.Price)
}

func (suite *TestSuite) TestTask_FinalizeTxs_Failed() {
	name := "TestTask_FinalizeTxs_Failed: invalid block range"
	suite.T().Log(name)

	// mock DB
	var (
		toBlock uint64 = 12379659
	)

	// mock DB
	db.Exec(`INSERT INTO public.pools (id, address, start_block, token0, token0_decimals, token1, token1_decimals, created_at) VALUES(1, '0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640', 12379656, '0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48', 6, '0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2', 18, 1712913866);`)
	db.Exec(`INSERT INTO public.block_cursors (id, pool_id, "type", block_number, extra, created_at, updated_at) VALUES(2, 1, 'finalizer', 12379656, '{"createdAtFinalizedBlock":12379581}'::json, 1712913868, 1712913962);
	INSERT INTO public.block_cursors (id, pool_id, "type", block_number, extra, created_at, updated_at) VALUES(1, 1, 'scanner', 12379660, NULL, 1712913868, 1712913963);`)

	db.Exec(`INSERT INTO public.txs (id, pool_id, tx_hash, block_number, block_time, gas, gas_price, receipt, is_finalized, created_at) VALUES(97, 1, '0x72e52be2d7f03b0842f4cc93e0fc798f340f71039bde322be821be7c8e8e677f', 12379657, 1712913887, 248711, 24859996161, '{}', false, 1712913893);`)
	db.Exec(`INSERT INTO public.swap_events (id, pool_id, tx_id, amount0, amount1, price, created_at) VALUES(96, 1, 97, '1160580000', '-329388970022326866', '3523.433101', 1712913893);`)

	// mock call
	ctrl := gomock.NewController(suite.T())
	defer ctrl.Finish()

	mockEthClient := eth.NewMockIClient(ctrl)
	mockEthClient.EXPECT().
		GetLatestBlockHeader(gomock.Any()).
		Return(&types.Header{Number: big.NewInt(int64(toBlock)), Time: uint64(time.Now().Unix())}, nil).
		Times(1)

	taskSvc.SetEthClient(mockEthClient)

	// handle
	payload := valueobject.TaskFinalizeTxsPayload{
		PoolID:         1,
		PoolAddress:    "0x88e6a0c2ddd26feeb64f039a2c41296fcb3f5640",
		Token0Decimals: 6,
		Token1Decimals: 18,
	}

	err := taskSvc.FinalizeTxs(context.Background(), payload)
	suite.ErrorIs(err, tasksvc.ErrInvalidBlockRange)

	// assert DB
	var cursor entity.BlockCursor
	db.First(&cursor)
	suite.EqualValues(12379660, cursor.BlockNumber)

	var tx entity.Tx
	db.First(&tx)
	suite.False(tx.IsFinalized)

	var event entity.SwapEvent
	err = db.First(&event).Error
	suite.Nil(err)
}
