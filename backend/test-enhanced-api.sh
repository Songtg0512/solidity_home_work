#!/bin/bash

# NFT 拍卖增强 API 测试脚本

BASE_URL="http://localhost:8080"

echo "=========================================="
echo "NFT 拍卖增强 API 测试"
echo "=========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 1. 测试增强的拍卖列表 API
echo -e "${YELLOW}1. 测试增强的拍卖列表 API${NC}"
echo -e "${BLUE}按最高出价排序：${NC}"
curl -s "$BASE_URL/api/auctions?status=active&sort_by=highest_bid&order=desc&page=1&page_size=5" | jq '.auctions[] | {auction_id, highest_bid, bid_count}'
echo ""

echo -e "${BLUE}按出价次数排序：${NC}"
curl -s "$BASE_URL/api/auctions?sort_by=bid_count&order=desc&page=1&page_size=5" | jq '.auctions[] | {auction_id, bid_count}'
echo ""

echo -e "${BLUE}按分类过滤：${NC}"
curl -s "$BASE_URL/api/auctions?category=art&page=1&page_size=5" | jq '.total, .auctions[].category' 
echo ""
echo ""

# 2. 测试增强统计信息（含 TVL）
echo -e "${YELLOW}2. 测试增强统计信息（TVL）${NC}"
curl -s "$BASE_URL/api/stats/enhanced" | jq
echo ""
echo ""

# 3. 测试获取钱包 NFT 列表
echo -e "${YELLOW}3. 测试获取钱包 NFT 列表${NC}"
read -p "请输入钱包地址 (留空跳过): " wallet_address
if [[ ! -z "$wallet_address" ]]; then
    echo -e "${BLUE}查询钱包 $wallet_address 的 NFT...${NC}"
    curl -s "$BASE_URL/api/wallet/$wallet_address/nfts" | jq '.ownedNfts[0:3] | .[] | {title, contract: .contract.address}'
    echo ""
fi
echo ""

# 4. 测试获取 NFT 地板价
echo -e "${YELLOW}4. 测试获取 NFT 地板价${NC}"
read -p "请输入 NFT 合约地址 (默认 BAYC): " nft_contract
nft_contract=${nft_contract:-0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D}
echo -e "${BLUE}查询合约 $nft_contract 的地板价...${NC}"
curl -s "$BASE_URL/api/nft/$nft_contract/floor-price" | jq
echo ""
echo ""

# 5. 测试获取 NFT 元数据
echo -e "${YELLOW}5. 测试获取 NFT 元数据${NC}"
read -p "请输入 NFT 合约地址 (默认 BAYC): " meta_contract
meta_contract=${meta_contract:-0xBC4CA0EdA7647A8aB7C2061c2E118A18a936f13D}
read -p "请输入 Token ID (默认 1): " token_id
token_id=${token_id:-1}
echo -e "${BLUE}查询 NFT 元数据...${NC}"
curl -s "$BASE_URL/api/nft/$meta_contract/$token_id/metadata" | jq '{name, description, image}'
echo ""
echo ""

# 6. 测试不同排序方式
echo -e "${YELLOW}6. 测试不同排序方式${NC}"

echo -e "${BLUE}按开始时间降序（最新的）：${NC}"
curl -s "$BASE_URL/api/auctions?sort_by=start_time&order=desc&page=1&page_size=3" | jq '.auctions[] | {auction_id, start_time}'
echo ""

echo -e "${BLUE}按起始价格升序（最便宜的）：${NC}"
curl -s "$BASE_URL/api/auctions?sort_by=start_price&order=asc&page=1&page_size=3" | jq '.auctions[] | {auction_id, start_price}'
echo ""

echo -e "${BLUE}按最高出价降序（最贵的）：${NC}"
curl -s "$BASE_URL/api/auctions?sort_by=highest_bid&order=desc&page=1&page_size=3" | jq '.auctions[] | {auction_id, highest_bid}'
echo ""
echo ""

# 7. 测试组合过滤
echo -e "${YELLOW}7. 测试组合过滤${NC}"
echo -e "${BLUE}活跃 + 按出价次数排序：${NC}"
curl -s "$BASE_URL/api/auctions?status=active&sort_by=bid_count&order=desc&page=1&page_size=3" | jq '.auctions[] | {auction_id, bid_count, ended}'
echo ""
echo ""

# 8. 性能测试
echo -e "${YELLOW}8. 性能测试${NC}"
echo -e "${BLUE}测试响应时间...${NC}"

# 测试拍卖列表
start_time=$(date +%s%N)
curl -s "$BASE_URL/api/auctions?page=1&page_size=10" > /dev/null
end_time=$(date +%s%N)
elapsed=$((($end_time - $start_time) / 1000000))
echo "拍卖列表 API: ${elapsed}ms"

# 测试统计信息
start_time=$(date +%s%N)
curl -s "$BASE_URL/api/stats/enhanced" > /dev/null
end_time=$(date +%s%N)
elapsed=$((($end_time - $start_time) / 1000000))
echo "统计信息 API: ${elapsed}ms"

echo ""
echo -e "${GREEN}=========================================="
echo "测试完成！"
echo "==========================================${NC}"
