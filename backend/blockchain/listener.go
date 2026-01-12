package blockchain

import (
	"auction-backend/config"
	"auction-backend/database"
	"auction-backend/models"
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// NftAuction 合约 ABI (只包含事件部分)
const NftAuctionABI = `[
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "internalType": "uint256", "name": "auctionId", "type": "uint256"},
			{"indexed": true, "internalType": "address", "name": "seller", "type": "address"},
			{"indexed": true, "internalType": "address", "name": "nftContract", "type": "address"},
			{"indexed": false, "internalType": "uint256", "name": "tokenId", "type": "uint256"},
			{"indexed": false, "internalType": "uint256", "name": "startPrice", "type": "uint256"},
			{"indexed": false, "internalType": "uint256", "name": "duration", "type": "uint256"},
			{"indexed": false, "internalType": "uint256", "name": "startTime", "type": "uint256"}
		],
		"name": "AuctionCreated",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "internalType": "uint256", "name": "auctionId", "type": "uint256"},
			{"indexed": true, "internalType": "address", "name": "bidder", "type": "address"},
			{"indexed": false, "internalType": "uint256", "name": "amount", "type": "uint256"},
			{"indexed": false, "internalType": "address", "name": "tokenAddress", "type": "address"},
			{"indexed": false, "internalType": "uint256", "name": "timestamp", "type": "uint256"}
		],
		"name": "BidPlaced",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "internalType": "uint256", "name": "auctionId", "type": "uint256"},
			{"indexed": true, "internalType": "address", "name": "winner", "type": "address"},
			{"indexed": false, "internalType": "uint256", "name": "finalPrice", "type": "uint256"},
			{"indexed": false, "internalType": "address", "name": "tokenAddress", "type": "address"},
			{"indexed": false, "internalType": "uint256", "name": "timestamp", "type": "uint256"}
		],
		"name": "AuctionEnded",
		"type": "event"
	}
]`

type EventListener struct {
	client          *ethclient.Client
	contractAddress common.Address
	contractABI     abi.ABI
}

// NewEventListener 创建事件监听器
func NewEventListener() (*EventListener, error) {
	client, err := ethclient.Dial(config.AppConfig.ETHRPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ethereum client: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(NftAuctionABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	return &EventListener{
		client:          client,
		contractAddress: common.HexToAddress(config.AppConfig.ContractAddress),
		contractABI:     contractABI,
	}, nil
}

// StartListening 开始监听事件
func (el *EventListener) StartListening(ctx context.Context) error {
	log.Println("Starting event listener...")

	// 获取最新已处理的区块号
	startBlock := config.AppConfig.StartBlock
	
	// 订阅新区块
	query := ethereum.FilterQuery{
		Addresses: []common.Address{el.contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := el.client.SubscribeFilterLogs(ctx, query, logs)
	if err != nil {
		// 如果订阅失败，使用轮询方式
		log.Println("Subscribe failed, using polling mode...")
		return el.pollLogs(ctx, startBlock)
	}

	log.Println("Event listener started successfully")

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Subscription error: %v\n", err)
			return err
		case vLog := <-logs:
			el.handleLog(vLog)
		case <-ctx.Done():
			log.Println("Event listener stopped")
			return nil
		}
	}
}

// pollLogs 轮询方式获取日志
func (el *EventListener) pollLogs(ctx context.Context, fromBlock uint64) error {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		Addresses: []common.Address{el.contractAddress},
	}

	logs, err := el.client.FilterLogs(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
	}

	for _, vLog := range logs {
		el.handleLog(vLog)
	}

	return nil
}

// handleLog 处理单个日志
func (el *EventListener) handleLog(vLog types.Log) {
	eventName := ""
	for name, event := range el.contractABI.Events {
		if event.ID == vLog.Topics[0] {
			eventName = name
			break
		}
	}

	switch eventName {
	case "AuctionCreated":
		el.handleAuctionCreated(vLog)
	case "BidPlaced":
		el.handleBidPlaced(vLog)
	case "AuctionEnded":
		el.handleAuctionEnded(vLog)
	default:
		log.Printf("Unknown event: %s\n", vLog.Topics[0].Hex())
	}
}

// handleAuctionCreated 处理拍卖创建事件
func (el *EventListener) handleAuctionCreated(vLog types.Log) {
	type AuctionCreatedEvent struct {
		AuctionId   *big.Int
		Seller      common.Address
		NftContract common.Address
		TokenId     *big.Int
		StartPrice  *big.Int
		Duration    *big.Int
		StartTime   *big.Int
	}

	var event AuctionCreatedEvent
	err := el.contractABI.UnpackIntoInterface(&event, "AuctionCreated", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack AuctionCreated event: %v\n", err)
		return
	}

	// 从 Topics 中提取 indexed 参数
	if len(vLog.Topics) >= 4 {
		event.AuctionId = new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		event.Seller = common.BytesToAddress(vLog.Topics[2].Bytes())
		event.NftContract = common.BytesToAddress(vLog.Topics[3].Bytes())
	}

	auction := models.Auction{
		AuctionID:   uint(event.AuctionId.Uint64()),
		Seller:      strings.ToLower(event.Seller.Hex()),
		NFTContract: strings.ToLower(event.NftContract.Hex()),
		TokenID:     event.TokenId.String(),
		StartPrice:  event.StartPrice.String(),
		Duration:    event.Duration.Uint64(),
		StartTime:   event.StartTime.Uint64(),
		Ended:       false,
	}

	db := database.GetDB()
	if err := db.Create(&auction).Error; err != nil {
		log.Printf("Failed to save auction: %v\n", err)
		return
	}

	log.Printf("Auction created: ID=%d, Seller=%s\n", auction.AuctionID, auction.Seller)
}

// handleBidPlaced 处理出价事件
func (el *EventListener) handleBidPlaced(vLog types.Log) {
	type BidPlacedEvent struct {
		AuctionId    *big.Int
		Bidder       common.Address
		Amount       *big.Int
		TokenAddress common.Address
		Timestamp    *big.Int
	}

	var event BidPlacedEvent
	err := el.contractABI.UnpackIntoInterface(&event, "BidPlaced", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack BidPlaced event: %v\n", err)
		return
	}

	// 从 Topics 中提取 indexed 参数
	if len(vLog.Topics) >= 3 {
		event.AuctionId = new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		event.Bidder = common.BytesToAddress(vLog.Topics[2].Bytes())
	}

	bid := models.Bid{
		AuctionID:    uint(event.AuctionId.Uint64()),
		Bidder:       strings.ToLower(event.Bidder.Hex()),
		Amount:       event.Amount.String(),
		TokenAddress: strings.ToLower(event.TokenAddress.Hex()),
		TxHash:       vLog.TxHash.Hex(),
		BlockNumber:  vLog.BlockNumber,
		Timestamp:    event.Timestamp.Uint64(),
	}

	db := database.GetDB()
	// 保存出价记录
	if err := db.Create(&bid).Error; err != nil {
		log.Printf("Failed to save bid: %v\n", err)
		return
	}

	// 更新拍卖的最高出价信息
	var auction models.Auction
	if err := db.Where("auction_id = ?", bid.AuctionID).First(&auction).Error; err != nil {
		log.Printf("Failed to find auction: %v\n", err)
		return
	}

	auction.HighestBidder = bid.Bidder
	auction.HighestBid = bid.Amount
	auction.TokenAddress = bid.TokenAddress
	auction.BidCount++  // 增加出价次数

	if err := db.Save(&auction).Error; err != nil {
		log.Printf("Failed to update auction: %v\n", err)
		return
	}

	log.Printf("Bid placed: AuctionID=%d, Bidder=%s, Amount=%s\n", bid.AuctionID, bid.Bidder, bid.Amount)
}

// handleAuctionEnded 处理拍卖结束事件
func (el *EventListener) handleAuctionEnded(vLog types.Log) {
	type AuctionEndedEvent struct {
		AuctionId    *big.Int
		Winner       common.Address
		FinalPrice   *big.Int
		TokenAddress common.Address
		Timestamp    *big.Int
	}

	var event AuctionEndedEvent
	err := el.contractABI.UnpackIntoInterface(&event, "AuctionEnded", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack AuctionEnded event: %v\n", err)
		return
	}

	// 从 Topics 中提取 indexed 参数
	if len(vLog.Topics) >= 3 {
		event.AuctionId = new(big.Int).SetBytes(vLog.Topics[1].Bytes())
		event.Winner = common.BytesToAddress(vLog.Topics[2].Bytes())
	}

	db := database.GetDB()
	var auction models.Auction
	if err := db.Where("auction_id = ?", event.AuctionId.Uint64()).First(&auction).Error; err != nil {
		log.Printf("Failed to find auction: %v\n", err)
		return
	}

	endTime := event.Timestamp.Uint64()
	auction.Ended = true
	auction.EndTime = &endTime
	auction.HighestBidder = strings.ToLower(event.Winner.Hex())
	auction.HighestBid = event.FinalPrice.String()
	auction.TokenAddress = strings.ToLower(event.TokenAddress.Hex())

	if err := db.Save(&auction).Error; err != nil {
		log.Printf("Failed to update auction: %v\n", err)
		return
	}

	log.Printf("Auction ended: ID=%d, Winner=%s, FinalPrice=%s\n", 
		auction.AuctionID, auction.HighestBidder, auction.HighestBid)
}
