# NFT 拍卖系统 API 文档

## 基础信息

- **Base URL**: `http://localhost:8080`
- **版本**: v1.0
- **认证**: 暂无（开发版本）

---

## 📋 API 端点列表

### 1. 拍卖相关 API

#### 1.1 获取拍卖列表（增强版）

```
GET /api/auctions
```

**查询参数：**

| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| page | int | 否 | 页码，默认1 | 1 |
| page_size | int | 否 | 每页数量，默认10，最大100 | 20 |
| status | string | 否 | 拍卖状态：`active`/`ended`/`all` | active |
| seller | string | 否 | 卖家地址过滤 | 0x123... |
| nft_contract | string | 否 | NFT合约地址过滤 | 0xabc... |
| category | string | 否 | 分类过滤 | art |
| sort_by | string | 否 | 排序字段：`start_time`/`highest_bid`/`bid_count`/`start_price` | highest_bid |
| order | string | 否 | 排序方向：`asc`/`desc`，默认`desc` | desc |

**请求示例：**
```bash
curl "http://localhost:8080/api/auctions?status=active&sort_by=highest_bid&order=desc&page=1&page_size=10"
```

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
      "bid_count": 5,
      "category": "art",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ]
}
```

#### 1.2 获取拍卖详情

```
GET /api/auctions/:id
```

**路径参数：**
- `id`: 拍卖ID

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
  "bid_count": 5,
  "category": "art",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

#### 1.3 获取拍卖出价历史

```
GET /api/auctions/:id/bids
```

**查询参数：**
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

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

---

### 2. NFT 相关 API

#### 2.1 获取钱包拥有的 NFT 列表

```
GET /api/wallet/:address/nfts
```

**路径参数：**
- `address`: 钱包地址

**查询参数：**
- `page_key`: Alchemy 分页键（可选）

**请求示例：**
```bash
curl "http://localhost:8080/api/wallet/0x1234.../nfts"
```

**响应示例：**
```json
{
  "ownedNfts": [
    {
      "contract": {
        "address": "0x5678..."
      },
      "id": "1",
      "title": "Bored Ape #1",
      "description": "A bored ape",
      "tokenUri": {
        "gateway": "https://..."
      },
      "media": [
        {
          "gateway": "https://image.url"
        }
      ],
      "metadata": {
        "name": "Bored Ape #1",
        "description": "...",
        "image": "https://...",
        "attributes": [...]
      }
    }
  ],
  "totalCount": 100,
  "pageKey": "next_page_key"
}
```

#### 2.2 获取 NFT 集合地板价

```
GET /api/nft/:contract/floor-price
```

**路径参数：**
- `contract`: NFT 合约地址

**请求示例：**
```bash
curl "http://localhost:8080/api/nft/0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D/floor-price"
```

**响应示例：**
```json
{
  "contract": "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d",
  "floor_price": "34.5",
  "volume_24h": "1234.56",
  "last_updated": "2024-01-01T12:00:00Z",
  "source": "opensea"
}
```

#### 2.3 获取 NFT 元数据

```
GET /api/nft/:contract/:token_id/metadata
```

**路径参数：**
- `contract`: NFT 合约地址
- `token_id`: Token ID

**请求示例：**
```bash
curl "http://localhost:8080/api/nft/0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D/1/metadata"
```

**响应示例：**
```json
{
  "id": 1,
  "contract": "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d",
  "token_id": "1",
  "name": "Bored Ape #1",
  "description": "A bored ape yacht club member",
  "image": "https://ipfs.io/ipfs/...",
  "attributes": "[{\"trait_type\":\"Background\",\"value\":\"Blue\"}]",
  "owner": "0x1234...",
  "floor_price": "34.5",
  "last_sync": "2024-01-01T12:00:00Z",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

---

### 3. 出价相关 API

#### 3.1 获取用户出价记录

```
GET /api/bids?bidder=0x...
```

**查询参数：**
- `bidder`: 出价者地址（必填）
- `page`: 页码（默认1）
- `page_size`: 每页数量（默认10）

**响应示例：**
```json
{
  "total": 25,
  "page": 1,
  "page_size": 10,
  "bids": [...]
}
```

#### 3.2 参与出价

```
POST /api/auctions/:id/bid
```

**请求体：**
```json
{
  "amount": "2000000000000000000",
  "token_address": "0x0000000000000000000000000000000000000000",
  "private_key": "your_private_key"
}
```

**响应示例：**
```json
{
  "message": "Bid placed successfully",
  "tx_hash": "0x1234..."
}
```

#### 3.3 结束拍卖

```
POST /api/auctions/:id/end
```

**请求体：**
```json
{
  "private_key": "your_private_key"
}
```

**响应示例：**
```json
{
  "message": "Auction ended successfully",
  "tx_hash": "0x5678..."
}
```

---

### 4. 统计信息 API

#### 4.1 获取基本统计信息

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

#### 4.2 获取增强统计信息（包含 TVL）

```
GET /api/stats/enhanced
```

**响应示例：**
```json
{
  "total_auctions": 100,
  "active_auctions": 30,
  "ended_auctions": 70,
  "total_bids": 500,
  "tvl": "50000000000000000000",
  "total_volume": "1000000000000000000000"
}
```

**字段说明：**
- `tvl`: Total Value Locked，所有活跃拍卖的最高出价总和（wei）
- `total_volume`: 所有出价的总和（wei）

---

### 5. 区块链交互 API

#### 5.1 从合约读取拍卖信息

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

---

## 📊 使用场景

### 场景 1: 显示首页拍卖列表

```bash
# 获取活跃拍卖，按出价从高到低排序
curl "http://localhost:8080/api/auctions?status=active&sort_by=highest_bid&order=desc&page=1&page_size=12"
```

### 场景 2: 显示拍卖详情页

```bash
# 1. 获取拍卖基本信息
curl "http://localhost:8080/api/auctions/0"

# 2. 获取 NFT 元数据（图片、名称等）
curl "http://localhost:8080/api/nft/0x5678.../1/metadata"

# 3. 获取出价历史
curl "http://localhost:8080/api/auctions/0/bids"

# 4. 获取集合地板价
curl "http://localhost:8080/api/nft/0x5678.../floor-price"
```

### 场景 3: 用户个人主页

```bash
# 1. 获取用户拥有的所有 NFT
curl "http://localhost:8080/api/wallet/0x1234.../nfts"

# 2. 获取用户的出价记录
curl "http://localhost:8080/api/bids?bidder=0x1234..."
```

### 场景 4: Dashboard 统计页面

```bash
# 获取完整统计信息
curl "http://localhost:8080/api/stats/enhanced"
```

---

## 🔧 错误处理

所有 API 错误返回格式：

```json
{
  "error": "错误描述信息"
}
```

**常见错误码：**
- `400`: 请求参数错误
- `404`: 资源不存在
- `500`: 服务器内部错误

---

## 💡 最佳实践

1. **分页查询**：
   - 建议 page_size 不超过 100
   - 使用缓存减少数据库压力

2. **地板价查询**：
   - 系统会缓存 1 小时
   - 频繁查询请使用缓存结果

3. **NFT 元数据**：
   - 系统会缓存 24 小时
   - 支持手动刷新

4. **私钥安全**：
   - ⚠️ 出价和结束拍卖接口仅用于开发测试
   - 生产环境应使用前端钱包直接与合约交互

---

## 📈 性能优化建议

1. **使用排序和过滤**：
   ```bash
   # 按最高出价排序，只查询活跃拍卖
   /api/auctions?status=active&sort_by=highest_bid&order=desc
   ```

2. **合理设置分页大小**：
   ```bash
   # 移动端建议每页 10 条
   /api/auctions?page_size=10
   
   # 桌面端可以每页 20-50 条
   /api/auctions?page_size=20
   ```

3. **利用缓存**：
   - 地板价：1小时缓存
   - NFT元数据：24小时缓存
   - 减少对第三方 API 的调用

---

## 🔗 相关链接

- [Alchemy API 文档](https://docs.alchemy.com/)
- [OpenSea API 文档](https://docs.opensea.io/)
- [项目 GitHub](https://github.com/...)
