package routes

import (
	"auction-backend/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", handlers.HealthCheck)

	// API 路由组
	api := r.Group("/api")
	{
		// 拍卖相关
		api.GET("/auctions", handlers.GetAuctionList)              // 获取拍卖列表（支持排序和分类）
		api.GET("/auctions/:id", handlers.GetAuctionDetail)        // 获取拍卖详情
		api.GET("/auctions/:id/bids", handlers.GetAuctionBids)     // 获取拍卖的出价历史
		api.GET("/auctions/:id/contract", handlers.GetContractAuctionInfo) // 从合约读取拍卖信息
		api.POST("/auctions/:id/bid", handlers.PlaceBid)           // 参与出价
		api.POST("/auctions/:id/end", handlers.EndAuction)         // 结束拍卖

		// 出价相关
		api.GET("/bids", handlers.GetBidsByBidder)                 // 获取某个地址的出价记录

		// NFT 相关
		api.GET("/wallet/:address/nfts", handlers.GetWalletNFTs)  // 获取钱包拥有的 NFT
		api.GET("/nft/:contract/floor-price", handlers.GetNFTFloorPrice) // 获取地板价
		api.GET("/nft/:contract/:token_id/metadata", handlers.GetNFTMetadata) // 获取 NFT 元数据

		// 统计信息
		api.GET("/stats", handlers.GetStats)                       // 获取基本统计信息
		api.GET("/stats/enhanced", handlers.GetEnhancedStats)      // 获取增强统计信息（含 TVL）
	}
}
