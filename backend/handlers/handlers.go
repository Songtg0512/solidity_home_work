package handlers

import (
	"auction-backend/blockchain"
	"auction-backend/config"
	"auction-backend/database"
	"auction-backend/models"
	"auction-backend/services"
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

// AuctionListResponse 拍卖列表响应
type AuctionListResponse struct {
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
	Auctions []models.Auction `json:"auctions"`
}

// BidListResponse 出价列表响应
type BidListResponse struct {
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
	Bids     []models.Bid `json:"bids"`
}

// GetAuctionList 获取拍卖列表
// GET /api/auctions?page=1&page_size=10&status=active&seller=0x...&sort_by=price&order=desc&category=art
func GetAuctionList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")                       // active, ended, all
	seller := c.Query("seller")                       // 卖家地址
	nftContract := c.Query("nft_contract")            // NFT合约地址
	category := c.Query("category")                   // 分类
	sortBy := c.DefaultQuery("sort_by", "start_time") // 排序字段: start_time, highest_bid, bid_count
	order := c.DefaultQuery("order", "desc")          // 排序顺序: asc, desc

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	db := database.GetDB()
	query := db.Model(&models.Auction{})

	// 过滤条件
	if status == "active" {
		query = query.Where("ended = ?", false)
	} else if status == "ended" {
		query = query.Where("ended = ?", true)
	}

	if seller != "" {
		query = query.Where("seller = ?", seller)
	}

	if nftContract != "" {
		query = query.Where("nft_contract = ?", nftContract)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 排序
	var orderClause string
	switch sortBy {
	case "highest_bid":
		orderClause = "CAST(highest_bid AS UNSIGNED)"
	case "bid_count":
		orderClause = "bid_count"
	case "start_price":
		orderClause = "CAST(start_price AS UNSIGNED)"
	default:
		orderClause = "start_time"
	}

	if order == "asc" {
		orderClause += " ASC"
	} else {
		orderClause += " DESC"
	}

	// 分页查询
	var auctions []models.Auction
	offset := (page - 1) * pageSize
	if err := query.Order(orderClause).
		Offset(offset).
		Limit(pageSize).
		Find(&auctions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query auctions",
		})
		return
	}

	c.JSON(http.StatusOK, AuctionListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Auctions: auctions,
	})
}

// GetAuctionDetail 获取拍卖详情
// GET /api/auctions/:id
func GetAuctionDetail(c *gin.Context) {
	auctionID := c.Param("id")

	db := database.GetDB()
	var auction models.Auction
	if err := db.Where("auction_id = ?", auctionID).First(&auction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Auction not found",
		})
		return
	}

	c.JSON(http.StatusOK, auction)
}

// GetAuctionBids 获取拍卖的出价历史
// GET /api/auctions/:id/bids?page=1&page_size=10
func GetAuctionBids(c *gin.Context) {
	auctionID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	db := database.GetDB()
	query := db.Model(&models.Bid{}).Where("auction_id = ?", auctionID)

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var bids []models.Bid
	offset := (page - 1) * pageSize
	if err := query.Order("timestamp DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query bids",
		})
		return
	}

	c.JSON(http.StatusOK, BidListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Bids:     bids,
	})
}

// GetBidsByBidder 获取某个地址的所有出价记录
// GET /api/bids?bidder=0x...&page=1&page_size=10
func GetBidsByBidder(c *gin.Context) {
	bidder := c.Query("bidder")
	if bidder == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bidder address is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	db := database.GetDB()
	query := db.Model(&models.Bid{}).Where("bidder = ?", bidder)

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var bids []models.Bid
	offset := (page - 1) * pageSize
	if err := query.Order("timestamp DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to query bids",
		})
		return
	}

	c.JSON(http.StatusOK, BidListResponse{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Bids:     bids,
	})
}

// GetStats 获取统计信息
// GET /api/stats
func GetStats(c *gin.Context) {
	db := database.GetDB()

	var totalAuctions int64
	var activeAuctions int64
	var endedAuctions int64
	var totalBids int64

	db.Model(&models.Auction{}).Count(&totalAuctions)
	db.Model(&models.Auction{}).Where("ended = ?", false).Count(&activeAuctions)
	db.Model(&models.Auction{}).Where("ended = ?", true).Count(&endedAuctions)
	db.Model(&models.Bid{}).Count(&totalBids)

	c.JSON(http.StatusOK, gin.H{
		"total_auctions":  totalAuctions,
		"active_auctions": activeAuctions,
		"ended_auctions":  endedAuctions,
		"total_bids":      totalBids,
	})
}

// HealthCheck 健康检查
// GET /health
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// PlaceBidRequest 出价请求
type PlaceBidRequest struct {
	AuctionID    string `json:"auction_id" binding:"required"`  // 拍卖ID
	Amount       string `json:"amount" binding:"required"`      // 出价金额（wei）
	TokenAddress string `json:"token_address"`                  // 代币地址，空或0x0表示ETH
	PrivateKey   string `json:"private_key" binding:"required"` // 用户私钥
}

// PlaceBid 参与出价
// POST /api/auctions/:id/bid
func PlaceBid(c *gin.Context) {
	auctionID := c.Param("id")

	var req PlaceBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 验证拍卖是否存在且未结束
	db := database.GetDB()
	var auction models.Auction
	if err := db.Where("auction_id = ?", auctionID).First(&auction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Auction not found",
		})
		return
	}

	if auction.Ended {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Auction has already ended",
		})
		return
	}

	// 检查拍卖是否超时
	currentTime := uint64(time.Now().Unix())
	if currentTime > auction.StartTime+auction.Duration {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Auction has expired",
		})
		return
	}

	// 解析金额
	amount, ok := new(big.Int).SetString(req.Amount, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid amount",
		})
		return
	}

	// 解析拍卖ID
	auctionIDInt, ok := new(big.Int).SetString(auctionID, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid auction ID",
		})
		return
	}

	// 处理代币地址
	var tokenAddress common.Address
	if req.TokenAddress == "" || req.TokenAddress == "0x0" || req.TokenAddress == "0x0000000000000000000000000000000000000000" {
		tokenAddress = common.Address{} // ETH
	} else {
		tokenAddress = common.HexToAddress(req.TokenAddress)
	}

	// 创建合约服务
	contractService, err := blockchain.NewContractService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to blockchain: " + err.Error(),
		})
		return
	}

	// 调用合约
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := contractService.PlaceBid(ctx, blockchain.PlaceBidRequest{
		AuctionID:    auctionIDInt,
		Amount:       amount,
		TokenAddress: tokenAddress,
		PrivateKey:   req.PrivateKey,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to place bid: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Bid placed successfully",
		"tx_hash": tx.Hash().Hex(),
	})
}

// EndAuctionRequest 结束拍卖请求
type EndAuctionRequest struct {
	PrivateKey string `json:"private_key" binding:"required"` // 调用者私钥
}

// EndAuction 结束拍卖
// POST /api/auctions/:id/end
func EndAuction(c *gin.Context) {
	auctionID := c.Param("id")

	var req EndAuctionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// 验证拍卖是否存在
	db := database.GetDB()
	var auction models.Auction
	if err := db.Where("auction_id = ?", auctionID).First(&auction).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Auction not found",
		})
		return
	}

	if auction.Ended {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Auction has already ended",
		})
		return
	}

	// 检查拍卖是否可以结束
	currentTime := uint64(time.Now().Unix())
	if currentTime <= auction.StartTime+auction.Duration {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Auction has not expired yet",
		})
		return
	}

	// 解析拍卖ID
	auctionIDInt, ok := new(big.Int).SetString(auctionID, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid auction ID",
		})
		return
	}

	// 创建合约服务
	contractService, err := blockchain.NewContractService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to blockchain: " + err.Error(),
		})
		return
	}

	// 调用合约
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := contractService.EndAuction(ctx, blockchain.EndAuctionRequest{
		AuctionID:  auctionIDInt,
		PrivateKey: req.PrivateKey,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to end auction: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Auction ended successfully",
		"tx_hash": tx.Hash().Hex(),
	})
}

// GetContractAuctionInfo 从合约读取拍卖信息
// GET /api/auctions/:id/contract
func GetContractAuctionInfo(c *gin.Context) {
	auctionID := c.Param("id")

	// 解析拍卖ID
	auctionIDInt, ok := new(big.Int).SetString(auctionID, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid auction ID",
		})
		return
	}

	// 创建合约服务
	contractService, err := blockchain.NewContractService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to connect to blockchain: " + err.Error(),
		})
		return
	}

	// 从合约读取信息
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info, err := contractService.GetAuctionInfo(ctx, auctionIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get auction info: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, info)
}

// GetWalletNFTs 获取钱包地址拥有的所有 NFT
// GET /api/wallet/:address/nfts?page_key=xxx
func GetWalletNFTs(c *gin.Context) {
	address := c.Param("address")
	pageKey := c.Query("page_key")

	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wallet address is required",
		})
		return
	}

	// 使用 Alchemy API 查询
	alchemyService := services.NewAlchemyService()
	result, err := alchemyService.GetNFTsByOwner(address, pageKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch NFTs: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetNFTFloorPrice 获取 NFT 集合地板价
// GET /api/nft/:contract/floor-price
func GetNFTFloorPrice(c *gin.Context) {
	contract := c.Param("contract")

	if contract == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "contract address is required",
		})
		return
	}

	// 优先从数据库查询
	db := database.GetDB()
	var collection models.NFTCollection
	err := db.Where("contract = ?", strings.ToLower(contract)).First(&collection).Error

	// 如果数据库中存在且数据较新（小于1小时），直接返回
	if err == nil && time.Since(collection.LastSync) < time.Hour {
		c.JSON(http.StatusOK, gin.H{
			"contract":     collection.Contract,
			"floor_price":  collection.FloorPrice,
			"volume_24h":   collection.Volume24h,
			"last_updated": collection.LastSync,
			"source":       "database",
		})
		return
	}

	// 从 OpenSea 查询
	openseaService := services.NewOpenSeaService(config.AppConfig.OpenSeaAPIKey)
	floorPrice, err := openseaService.GetFloorPriceByContract(contract)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch floor price: " + err.Error(),
		})
		return
	}

	// 更新数据库
	if err == nil {
		collection.FloorPrice = fmt.Sprintf("%.18f", floorPrice)
		collection.LastSync = time.Now()
		db.Save(&collection)
	} else {
		collection = models.NFTCollection{
			Contract:   strings.ToLower(contract),
			FloorPrice: fmt.Sprintf("%.18f", floorPrice),
			LastSync:   time.Now(),
		}
		db.Create(&collection)
	}

	c.JSON(http.StatusOK, gin.H{
		"contract":     collection.Contract,
		"floor_price":  collection.FloorPrice,
		"last_updated": collection.LastSync,
		"source":       "opensea",
	})
}

// GetNFTMetadata 获取 NFT 元数据
// GET /api/nft/:contract/:token_id/metadata
func GetNFTMetadata(c *gin.Context) {
	contract := c.Param("contract")
	tokenID := c.Param("token_id")

	if contract == "" || tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "contract and token_id are required",
		})
		return
	}

	// 优先从数据库查询
	db := database.GetDB()
	var metadata models.NFTMetadata
	err := db.Where("contract = ? AND token_id = ?", strings.ToLower(contract), tokenID).First(&metadata).Error

	// 如果数据库中存在且数据较新（小于24小时），直接返回
	if err == nil && time.Since(metadata.LastSync) < 24*time.Hour {
		c.JSON(http.StatusOK, metadata)
		return
	}

	// 从 Alchemy 查询
	alchemyService := services.NewAlchemyService()
	newMetadata, err := alchemyService.GetNFTMetadata(contract, tokenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch metadata: " + err.Error(),
		})
		return
	}

	// 保存或更新数据库
	if err == nil {
		newMetadata.ID = metadata.ID
		db.Save(newMetadata)
	} else {
		db.Create(newMetadata)
	}

	c.JSON(http.StatusOK, newMetadata)
}

// GetEnhancedStats 获取增强的统计信息（包括 TVL）
// GET /api/stats/enhanced
func GetEnhancedStats(c *gin.Context) {
	db := database.GetDB()

	var totalAuctions int64
	var activeAuctions int64
	var endedAuctions int64
	var totalBids int64

	db.Model(&models.Auction{}).Count(&totalAuctions)
	db.Model(&models.Auction{}).Where("ended = ?", false).Count(&activeAuctions)
	db.Model(&models.Auction{}).Where("ended = ?", true).Count(&endedAuctions)
	db.Model(&models.Bid{}).Count(&totalBids)

	// 计算 TVL（所有活跃拍卖的最高出价总和）
	var activeAuctionsList []models.Auction
	db.Where("ended = ?", false).Find(&activeAuctionsList)

	tvl := big.NewInt(0)
	for _, auction := range activeAuctionsList {
		if auction.HighestBid != "" {
			amount, ok := new(big.Int).SetString(auction.HighestBid, 10)
			if ok {
				tvl.Add(tvl, amount)
			}
		}
	}

	// 计算总交易量（所有出价的总和）
	var allBids []models.Bid
	db.Find(&allBids)

	totalVolume := big.NewInt(0)
	for _, bid := range allBids {
		amount, ok := new(big.Int).SetString(bid.Amount, 10)
		if ok {
			totalVolume.Add(totalVolume, amount)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_auctions":  totalAuctions,
		"active_auctions": activeAuctions,
		"ended_auctions":  endedAuctions,
		"total_bids":      totalBids,
		"tvl":             tvl.String(),
		"total_volume":    totalVolume.String(),
	})
}
