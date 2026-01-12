package models

import (
	"time"
)

// Auction 拍卖表
type Auction struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	AuctionID     uint      `gorm:"uniqueIndex;not null" json:"auction_id"` // 链上拍卖ID
	Seller        string    `gorm:"size:42;not null;index" json:"seller"`
	NFTContract   string    `gorm:"size:42;not null;index" json:"nft_contract"`
	TokenID       string    `gorm:"size:78;not null" json:"token_id"`
	StartPrice    string    `gorm:"size:78;not null" json:"start_price"`
	Duration      uint64    `gorm:"not null" json:"duration"`
	StartTime     uint64    `gorm:"not null;index" json:"start_time"`
	Ended         bool      `gorm:"default:false;index" json:"ended"`
	HighestBidder string    `gorm:"size:42" json:"highest_bidder"`
	HighestBid    string    `gorm:"size:78" json:"highest_bid"`
	TokenAddress  string    `gorm:"size:42" json:"token_address"` // 出价代币地址，0x0为ETH
	EndTime       *uint64   `json:"end_time"`                     // 实际结束时间
	BidCount      int       `gorm:"default:0" json:"bid_count"`   // 出价次数
	Category      string    `gorm:"size:50;index" json:"category"` // 分类
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	
	// 关联的 NFT 元数据（非数据库字段）
	NFTMetadata *NFTMetadata `gorm:"-" json:"nft_metadata,omitempty"`
}

// Bid 出价记录表
type Bid struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	AuctionID    uint      `gorm:"not null;index" json:"auction_id"` // 链上拍卖ID
	Bidder       string    `gorm:"size:42;not null;index" json:"bidder"`
	Amount       string    `gorm:"size:78;not null" json:"amount"`
	TokenAddress string    `gorm:"size:42;not null" json:"token_address"`
	TxHash       string    `gorm:"size:66;uniqueIndex" json:"tx_hash"`
	BlockNumber  uint64    `gorm:"not null;index" json:"block_number"`
	Timestamp    uint64    `gorm:"not null;index" json:"timestamp"`
	CreatedAt    time.Time `json:"created_at"`
}

// NFTMetadata NFT 元数据表
type NFTMetadata struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Contract    string    `gorm:"size:42;not null;index" json:"contract"`
	TokenID     string    `gorm:"size:78;not null;index" json:"token_id"`
	Name        string    `gorm:"size:255" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Image       string    `gorm:"size:512" json:"image"`
	Attributes  string    `gorm:"type:text" json:"attributes"` // JSON 字符串
	Owner       string    `gorm:"size:42;index" json:"owner"`
	FloorPrice  string    `gorm:"size:78" json:"floor_price"` // 地板价
	LastSync    time.Time `json:"last_sync"`                  // 最后同步时间
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NFTCollection NFT 集合信息表
type NFTCollection struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Contract    string    `gorm:"size:42;uniqueIndex;not null" json:"contract"`
	Name        string    `gorm:"size:255" json:"name"`
	Symbol      string    `gorm:"size:50" json:"symbol"`
	TotalSupply int64     `json:"total_supply"`
	FloorPrice  string    `gorm:"size:78" json:"floor_price"`
	Volume24h   string    `gorm:"size:78" json:"volume_24h"`
	Description string    `gorm:"type:text" json:"description"`
	Image       string    `gorm:"size:512" json:"image"`
	LastSync    time.Time `json:"last_sync"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Auction) TableName() string {
	return "auctions"
}

func (Bid) TableName() string {
	return "bids"
}

func (NFTMetadata) TableName() string {
	return "nft_metadata"
}

func (NFTCollection) TableName() string {
	return "nft_collections"
}
