# NFT 拍卖后台服务

基于 Go 语言开发的 NFT 拍卖后台服务，用于监听链上拍卖合约事件并提供 REST API 查询接口。

## 功能特性

- ✅ 实时监听区块链事件（拍卖创建、出价、拍卖结束）
- ✅ 自动同步链上数据到 MySQL 数据库
- ✅ 提供 RESTful API 查询接口
- ✅ 支持分页查询和条件过滤
- ✅ 支持 ETH 和 ERC20 代币出价

## 技术栈

- **Go 1.21+**
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **go-ethereum** - 以太坊客户端
- **MySQL 8.0+** - 数据库

## 项目结构

```
backend/
├── main.go                 # 程序入口
├── config/                 # 配置管理
│   └── config.go
├── database/               # 数据库连接
│   └── database.go
├── models/                 # 数据模型
│   └── models.go
├── blockchain/             # 区块链事件监听
│   └── listener.go
├── handlers/               # API 处理器
│   └── handlers.go
├── routes/                 # 路由定义
│   └── routes.go
├── schema.sql             # 数据库表结构
├── .env.example           # 环境变量示例
└── go.mod                 # Go 模块依赖
```

## 快速开始

### 1. 安装依赖

```bash
cd backend
go mod download
```

### 2. 配置数据库

创建 MySQL 数据库并导入表结构：

```bash
mysql -u root -p < schema.sql
```

### 3. 配置环境变量

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cp .env.example .env
```

编辑 `.env` 文件：

```env
# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=nft_auction

# 区块链配置
ETH_RPC_URL=https://sepolia.infura.io/v3/your_infura_key
CONTRACT_ADDRESS=0x...  # NftAuction 合约地址
START_BLOCK=0           # 开始监听的区块号

# 服务器配置
SERVER_PORT=8080
```

### 4. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 接口文档

### 1. 健康检查

```
GET /health
```

**响应示例：**
```json
{
  "status": "ok"
}
```

### 2. 获取拍卖列表

```
GET /api/auctions?page=1&page_size=10&status=active&seller=0x...&nft_contract=0x...
```

**查询参数：**
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 10，最大 100）
- `status`: 拍卖状态（`active`/`ended`/`all`）
- `seller`: 卖家地址（可选）
- `nft_contract`: NFT 合约地址（可选）

**响应示例：**
```json
{
  "total": 100,
  "page": 1,
  "page_size": 10,
  "auctions": [
    {
      "id": 1,
      "auction_id": 0,
      "seller": "0x1234...",
      "nft_contract": "0x5678...",
      "token_id": "1",
      "start_price": "1000000000000000000",
      "duration": 86400,
      "start_time": 1704067200,
      "ended": false,
      "highest_bidder": "0xabcd...",
      "highest_bid": "2000000000000000000",
      "token_address": "0x0000000000000000000000000000000000000000",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

### 3. 获取拍卖详情

```
GET /api/auctions/:id
```

**响应示例：**
```json
{
  "id": 1,
  "auction_id": 0,
  "seller": "0x1234...",
  "nft_contract": "0x5678...",
  "token_id": "1",
  "start_price": "1000000000000000000",
  "duration": 86400,
  "start_time": 1704067200,
  "ended": false,
  "highest_bidder": "0xabcd...",
  "highest_bid": "2000000000000000000",
  "token_address": "0x0000000000000000000000000000000000000000",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### 4. 获取拍卖出价历史

```
GET /api/auctions/:id/bids?page=1&page_size=10
```

**响应示例：**
```json
{
  "total": 50,
  "page": 1,
  "page_size": 10,
  "bids": [
    {
      "id": 1,
      "auction_id": 0,
      "bidder": "0xabcd...",
      "amount": "2000000000000000000",
      "token_address": "0x0000000000000000000000000000000000000000",
      "tx_hash": "0x1234...",
      "block_number": 1000000,
      "timestamp": 1704070800,
      "created_at": "2024-01-01T01:00:00Z"
    }
  ]
}
```

### 5. 获取用户出价记录

```
GET /api/bids?bidder=0x...&page=1&page_size=10
```

**查询参数：**
- `bidder`: 出价者地址（必填）
- `page`: 页码（默认 1）
- `page_size`: 每页数量（默认 10，最大 100）

**响应格式同上**

### 6. 获取统计信息

```
GET /api/stats
```

**响应示例：**
```json
{
  "total_auctions": 100,
  "active_auctions": 30,
  "ended_auctions": 70,
  "total_bids": 500
}
```

### 7. 参与出价

```
POST /api/auctions/:id/bid
```

**请求体：**
```json
{
  "auction_id": "0",
  "amount": "2000000000000000000",
  "token_address": "0x0000000000000000000000000000000000000000",
  "private_key": "your_private_key_without_0x_prefix"
}
```

**字段说明：**
- `auction_id`: 拍卖ID（可选，优先使用URL参数）
- `amount`: 出价金额（wei 单位）
- `token_address`: 代币地址，空或 `0x0` 表示 ETH
- `private_key`: 用户私钥（不带 0x 前缀）

**响应示例：**
```json
{
  "message": "Bid placed successfully",
  "tx_hash": "0x1234567890abcdef..."
}
```

**错误响应：**
```json
{
  "error": "Auction has already ended"
}
```

### 8. 结束拍卖

```
POST /api/auctions/:id/end
```

**请求体：**
```json
{
  "private_key": "your_private_key_without_0x_prefix"
}
```

**字段说明：**
- `private_key`: 调用者私钥（不带 0x 前缀）

**响应示例：**
```json
{
  "message": "Auction ended successfully",
  "tx_hash": "0x1234567890abcdef..."
}
```

**错误响应：**
```json
{
  "error": "Auction has not expired yet"
}
```

### 9. 从合约读取拍卖信息

```
GET /api/auctions/:id/contract
```

**响应示例：**
```json
{
  "seller": "0x1234...",
  "duration": "86400",
  "start_price": "1000000000000000000",
  "start_time": "1704067200",
  "ended": false,
  "highest_bidder": "0xabcd...",
  "highest_bid": "2000000000000000000",
  "nft_contract": "0x5678...",
  "token_id": "1",
  "token_address": "0x0000000000000000000000000000000000000000"
}
```

## 数据库表结构

### auctions（拍卖表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| auction_id | BIGINT | 链上拍卖ID |
| seller | VARCHAR(42) | 卖家地址 |
| nft_contract | VARCHAR(42) | NFT合约地址 |
| token_id | VARCHAR(78) | NFT TokenID |
| start_price | VARCHAR(78) | 起始价格 |
| duration | BIGINT | 拍卖持续时间（秒） |
| start_time | BIGINT | 开始时间（Unix时间戳） |
| ended | BOOLEAN | 是否已结束 |
| highest_bidder | VARCHAR(42) | 最高出价者 |
| highest_bid | VARCHAR(78) | 最高出价 |
| token_address | VARCHAR(42) | 出价代币地址 |
| end_time | BIGINT | 实际结束时间 |

### bids（出价记录表）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | BIGINT | 主键 |
| auction_id | BIGINT | 链上拍卖ID |
| bidder | VARCHAR(42) | 出价者地址 |
| amount | VARCHAR(78) | 出价金额 |
| token_address | VARCHAR(42) | 出价代币地址 |
| tx_hash | VARCHAR(66) | 交易哈希 |
| block_number | BIGINT | 区块号 |
| timestamp | BIGINT | 时间戳 |

## 合约事件

后台服务监听以下合约事件：

### AuctionCreated
```solidity
event AuctionCreated(
    uint256 indexed auctionId,
    address indexed seller,
    address indexed nftContract,
    uint256 tokenId,
    uint256 startPrice,
    uint256 duration,
    uint256 startTime
);
```

### BidPlaced
```solidity
event BidPlaced(
    uint256 indexed auctionId,
    address indexed bidder,
    uint256 amount,
    address tokenAddress,
    uint256 timestamp
);
```

### AuctionEnded
```solidity
event AuctionEnded(
    uint256 indexed auctionId,
    address indexed winner,
    uint256 finalPrice,
    address tokenAddress,
    uint256 timestamp
);
```

## 部署说明

### 使用 Docker 部署

创建 `Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./main"]
```

构建并运行：

```bash
docker build -t auction-backend .
docker run -p 8080:8080 --env-file .env auction-backend
```

### 使用 systemd 部署

创建服务文件 `/etc/systemd/system/auction-backend.service`：

```ini
[Unit]
Description=NFT Auction Backend Service
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/auction-backend
ExecStart=/opt/auction-backend/main
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable auction-backend
sudo systemctl start auction-backend
```

## 注意事项

1. **区块同步**：首次启动时，需要设置 `START_BLOCK` 为合约部署的区块号，避免从创世区块开始扫描。

2. **RPC 限制**：公共 RPC 节点可能有请求限制，建议使用 Infura、Alchemy 等服务或自建节点。

3. **数据一致性**：建议定期备份数据库，防止数据丢失。

4. **性能优化**：
   - 使用数据库索引提高查询性能
   - 考虑使用 Redis 缓存热点数据
   - 使用连接池管理数据库连接

## 开发调试

启用 Gin 的调试模式：

```go
gin.SetMode(gin.DebugMode)
```

查看日志：

```bash
tail -f auction-backend.log
```

## License

MIT
