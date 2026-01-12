# API 使用示例

## 基础查询功能

### 1. 查询所有活跃的拍卖

```bash
curl http://localhost:8080/api/auctions?status=active&page=1&page_size=10
```

### 2. 查询某个拍卖的详情

```bash
curl http://localhost:8080/api/auctions/0
```

### 3. 查询某个拍卖的出价历史

```bash
curl http://localhost:8080/api/auctions/0/bids
```

### 4. 从区块链合约直接读取拍卖信息

```bash
curl http://localhost:8080/api/auctions/0/contract
```

### 5. 查询某个用户的所有出价记录

```bash
curl "http://localhost:8080/api/bids?bidder=0x1234567890123456789012345678901234567890"
```

## 交互功能（需要私钥）

⚠️ **安全警告**：在生产环境中，不应该将私钥发送到后端服务器。这些功能仅用于演示。实际应用中应该在前端使用 Web3/ethers.js 直接与合约交互。

### 1. 参与出价（使用 ETH）

```bash
curl -X POST http://localhost:8080/api/auctions/0/bid \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "2000000000000000000",
    "token_address": "0x0000000000000000000000000000000000000000",
    "private_key": "你的私钥（不带0x前缀）"
  }'
```

**响应示例：**
```json
{
  "message": "Bid placed successfully",
  "tx_hash": "0xabcdef..."
}
```

### 2. 参与出价（使用 ERC20 代币）

```bash
curl -X POST http://localhost:8080/api/auctions/0/bid \
  -H "Content-Type: application/json" \
  -d '{
    "amount": "1000000",
    "token_address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
    "private_key": "你的私钥（不带0x前缀）"
  }'
```

### 3. 结束拍卖

```bash
curl -X POST http://localhost:8080/api/auctions/0/end \
  -H "Content-Type: application/json" \
  -d '{
    "private_key": "你的私钥（不带0x前缀）"
  }'
```

**响应示例：**
```json
{
  "message": "Auction ended successfully",
  "tx_hash": "0x123456..."
}
```

## JavaScript/TypeScript 示例

### 使用 fetch API

```javascript
// 1. 查询拍卖列表
async function getAuctions() {
  const response = await fetch('http://localhost:8080/api/auctions?status=active');
  const data = await response.json();
  console.log('活跃拍卖:', data);
}

// 2. 参与出价（仅演示，生产环境不应这样做）
async function placeBid(auctionId, amount, privateKey) {
  const response = await fetch(`http://localhost:8080/api/auctions/${auctionId}/bid`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      amount: amount,
      token_address: '0x0000000000000000000000000000000000000000',
      private_key: privateKey
    })
  });
  
  const result = await response.json();
  console.log('出价结果:', result);
  return result;
}

// 3. 获取拍卖详情
async function getAuctionDetail(auctionId) {
  const response = await fetch(`http://localhost:8080/api/auctions/${auctionId}`);
  const data = await response.json();
  console.log('拍卖详情:', data);
}
```

### 使用 axios

```javascript
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

// 1. 获取拍卖列表
async function getAuctions(page = 1, pageSize = 10, status = 'active') {
  try {
    const response = await axios.get(`${API_BASE_URL}/auctions`, {
      params: { page, page_size: pageSize, status }
    });
    return response.data;
  } catch (error) {
    console.error('获取拍卖列表失败:', error);
    throw error;
  }
}

// 2. 获取出价历史
async function getBidHistory(auctionId, page = 1, pageSize = 10) {
  try {
    const response = await axios.get(`${API_BASE_URL}/auctions/${auctionId}/bids`, {
      params: { page, page_size: pageSize }
    });
    return response.data;
  } catch (error) {
    console.error('获取出价历史失败:', error);
    throw error;
  }
}

// 3. 获取统计信息
async function getStats() {
  try {
    const response = await axios.get(`${API_BASE_URL}/stats`);
    return response.data;
  } catch (error) {
    console.error('获取统计信息失败:', error);
    throw error;
  }
}
```

## Python 示例

```python
import requests
import json

API_BASE_URL = 'http://localhost:8080/api'

# 1. 获取拍卖列表
def get_auctions(page=1, page_size=10, status='active'):
    params = {
        'page': page,
        'page_size': page_size,
        'status': status
    }
    response = requests.get(f'{API_BASE_URL}/auctions', params=params)
    return response.json()

# 2. 获取拍卖详情
def get_auction_detail(auction_id):
    response = requests.get(f'{API_BASE_URL}/auctions/{auction_id}')
    return response.json()

# 3. 获取出价历史
def get_bid_history(auction_id, page=1, page_size=10):
    params = {
        'page': page,
        'page_size': page_size
    }
    response = requests.get(f'{API_BASE_URL}/auctions/{auction_id}/bids', params=params)
    return response.json()

# 4. 参与出价（仅演示）
def place_bid(auction_id, amount, token_address, private_key):
    data = {
        'amount': amount,
        'token_address': token_address,
        'private_key': private_key
    }
    response = requests.post(
        f'{API_BASE_URL}/auctions/{auction_id}/bid',
        json=data
    )
    return response.json()

# 5. 结束拍卖
def end_auction(auction_id, private_key):
    data = {
        'private_key': private_key
    }
    response = requests.post(
        f'{API_BASE_URL}/auctions/{auction_id}/end',
        json=data
    )
    return response.json()

# 使用示例
if __name__ == '__main__':
    # 获取活跃拍卖列表
    auctions = get_auctions(status='active')
    print('活跃拍卖数量:', auctions['total'])
    
    # 获取第一个拍卖的详情
    if auctions['auctions']:
        first_auction = auctions['auctions'][0]
        detail = get_auction_detail(first_auction['auction_id'])
        print('拍卖详情:', json.dumps(detail, indent=2))
        
        # 获取出价历史
        bids = get_bid_history(first_auction['auction_id'])
        print('出价数量:', bids['total'])
```

## Go 示例

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

const APIBaseURL = "http://localhost:8080/api"

type AuctionListResponse struct {
    Total    int64     `json:"total"`
    Page     int       `json:"page"`
    PageSize int       `json:"page_size"`
    Auctions []Auction `json:"auctions"`
}

type Auction struct {
    ID            uint   `json:"id"`
    AuctionID     uint   `json:"auction_id"`
    Seller        string `json:"seller"`
    NFTContract   string `json:"nft_contract"`
    TokenID       string `json:"token_id"`
    StartPrice    string `json:"start_price"`
    HighestBid    string `json:"highest_bid"`
    Ended         bool   `json:"ended"`
}

// 获取拍卖列表
func GetAuctions(status string, page, pageSize int) (*AuctionListResponse, error) {
    url := fmt.Sprintf("%s/auctions?status=%s&page=%d&page_size=%d", 
        APIBaseURL, status, page, pageSize)
    
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result AuctionListResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

// 参与出价
func PlaceBid(auctionID, amount, tokenAddress, privateKey string) (map[string]interface{}, error) {
    data := map[string]string{
        "amount":        amount,
        "token_address": tokenAddress,
        "private_key":   privateKey,
    }
    
    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    
    url := fmt.Sprintf("%s/auctions/%s/bid", APIBaseURL, auctionID)
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    
    return result, nil
}

func main() {
    // 获取活跃拍卖
    auctions, err := GetAuctions("active", 1, 10)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("Total active auctions: %d\n", auctions.Total)
    for _, auction := range auctions.Auctions {
        fmt.Printf("Auction #%d: %s\n", auction.AuctionID, auction.NFTContract)
    }
}
```

## 注意事项

1. **安全性**：
   - 🚨 **永远不要**在生产环境中将私钥发送到后端服务器
   - 实际应用应该使用 MetaMask 或其他钱包直接与合约交互
   - 这些私钥功能仅用于开发测试

2. **金额单位**：
   - 所有金额都使用 **wei** 单位
   - 1 ETH = 1,000,000,000,000,000,000 wei (10^18)
   - 可以使用 `ethers.utils.parseEther("1.0")` 转换

3. **地址格式**：
   - 所有地址都应该是完整的 42 字符格式（包括 0x 前缀）
   - ETH 使用 `0x0000000000000000000000000000000000000000`

4. **错误处理**：
   - 始终检查 HTTP 状态码
   - 解析 JSON 响应中的 `error` 字段
   - 处理网络超时和重试逻辑
