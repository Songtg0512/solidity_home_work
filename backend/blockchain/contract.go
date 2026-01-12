package blockchain

import (
	"auction-backend/config"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ContractService 合约服务
type ContractService struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
}

// NewContractService 创建合约服务实例
func NewContractService() (*ContractService, error) {
	client, err := ethclient.Dial(config.AppConfig.ETHRPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum client: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(NftAuctionABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	return &ContractService{
		client:          client,
		contractAddress: common.HexToAddress(config.AppConfig.ContractAddress),
		contractABI:     contractABI,
	}, nil
}

// PlaceBidRequest 出价请求
type PlaceBidRequest struct {
	AuctionID    *big.Int       // 拍卖ID
	Amount       *big.Int       // 出价金额
	TokenAddress common.Address // 代币地址，0x0 表示 ETH
	PrivateKey   string         // 用户私钥
}

// PlaceBid 参与出价
func (cs *ContractService) PlaceBid(ctx context.Context, req PlaceBidRequest) (*types.Transaction, error) {
	// 解析私钥
	privateKey, err := crypto.HexToECDSA(req.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// 获取公钥地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取 nonce
	nonce, err := cs.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// 获取 gas price
	gasPrice, err := cs.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// 获取 chain ID
	chainID, err := cs.client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// 创建交易选项
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(300000) // 设置 gas limit

	// 如果是 ETH 出价，设置 value
	if req.TokenAddress == (common.Address{}) {
		auth.Value = req.Amount
	} else {
		auth.Value = big.NewInt(0)
	}

	// 编码函数调用数据
	data, err := cs.contractABI.Pack("placeBid", req.AuctionID, req.Amount, req.TokenAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transaction data: %w", err)
	}

	// 创建交易
	tx := types.NewTransaction(
		nonce,
		cs.contractAddress,
		auth.Value,
		auth.GasLimit,
		gasPrice,
		data,
	)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 发送交易
	err = cs.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}

// EndAuctionRequest 结束拍卖请求
type EndAuctionRequest struct {
	AuctionID  *big.Int // 拍卖ID
	PrivateKey string   // 调用者私钥
}

// EndAuction 结束拍卖
func (cs *ContractService) EndAuction(ctx context.Context, req EndAuctionRequest) (*types.Transaction, error) {
	// 解析私钥
	privateKey, err := crypto.HexToECDSA(req.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// 获取公钥地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 获取 nonce
	nonce, err := cs.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	// 获取 gas price
	gasPrice, err := cs.client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}

	// 获取 chain ID
	chainID, err := cs.client.NetworkID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	// 编码函数调用数据
	data, err := cs.contractABI.Pack("endAuction", req.AuctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to pack transaction data: %w", err)
	}

	// 创建交易
	tx := types.NewTransaction(
		nonce,
		cs.contractAddress,
		big.NewInt(0),
		uint64(300000), // gas limit
		gasPrice,
		data,
	)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %w", err)
	}

	// 发送交易
	err = cs.client.SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send transaction: %w", err)
	}

	return signedTx, nil
}

// GetAuctionInfo 获取拍卖信息（从合约读取）
func (cs *ContractService) GetAuctionInfo(ctx context.Context, auctionID *big.Int) (map[string]interface{}, error) {
	// 编码函数调用
	data, err := cs.contractABI.Pack("auctions", auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to pack call data: %w", err)
	}

	// 调用合约
	result, err := cs.client.CallContract(ctx, map[string]interface{}{
		"to":   cs.contractAddress,
		"data": data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}

	// 解析结果
	var out struct {
		Seller        common.Address
		Duration      *big.Int
		StartPrice    *big.Int
		StartTime     *big.Int
		Ended         bool
		HighestBidder common.Address
		HighestBid    *big.Int
		NftContract   common.Address
		TokenId       *big.Int
		TokenAddress  common.Address
	}

	err = cs.contractABI.UnpackIntoInterface(&out, "auctions", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	return map[string]interface{}{
		"seller":         out.Seller.Hex(),
		"duration":       out.Duration.String(),
		"start_price":    out.StartPrice.String(),
		"start_time":     out.StartTime.String(),
		"ended":          out.Ended,
		"highest_bidder": out.HighestBidder.Hex(),
		"highest_bid":    out.HighestBid.String(),
		"nft_contract":   out.NftContract.Hex(),
		"token_id":       out.TokenId.String(),
		"token_address":  out.TokenAddress.Hex(),
	}, nil
}
