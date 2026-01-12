-- 创建数据库
CREATE DATABASE IF NOT EXISTS nft_auction CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE nft_auction;

-- 拍卖表
CREATE TABLE IF NOT EXISTS auctions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    auction_id BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '链上拍卖ID',
    seller VARCHAR(42) NOT NULL COMMENT '卖家地址',
    nft_contract VARCHAR(42) NOT NULL COMMENT 'NFT合约地址',
    token_id VARCHAR(78) NOT NULL COMMENT 'NFT TokenID',
    start_price VARCHAR(78) NOT NULL COMMENT '起始价格',
    duration BIGINT UNSIGNED NOT NULL COMMENT '拍卖持续时间(秒)',
    start_time BIGINT UNSIGNED NOT NULL COMMENT '开始时间(Unix时间戳)',
    ended BOOLEAN DEFAULT FALSE COMMENT '是否已结束',
    highest_bidder VARCHAR(42) COMMENT '最高出价者',
    highest_bid VARCHAR(78) COMMENT '最高出价',
    token_address VARCHAR(42) COMMENT '出价代币地址',
    end_time BIGINT UNSIGNED COMMENT '实际结束时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_seller (seller),
    INDEX idx_nft_contract (nft_contract),
    INDEX idx_start_time (start_time),
    INDEX idx_ended (ended)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='拍卖表';

-- 出价记录表
CREATE TABLE IF NOT EXISTS bids (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    auction_id BIGINT UNSIGNED NOT NULL COMMENT '链上拍卖ID',
    bidder VARCHAR(42) NOT NULL COMMENT '出价者地址',
    amount VARCHAR(78) NOT NULL COMMENT '出价金额',
    token_address VARCHAR(42) NOT NULL COMMENT '出价代币地址',
    tx_hash VARCHAR(66) NOT NULL UNIQUE COMMENT '交易哈希',
    block_number BIGINT UNSIGNED NOT NULL COMMENT '区块号',
    timestamp BIGINT UNSIGNED NOT NULL COMMENT '时间戳',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_auction_id (auction_id),
    INDEX idx_bidder (bidder),
    INDEX idx_block_number (block_number),
    INDEX idx_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='出价记录表';
