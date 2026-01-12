#!/bin/bash

# NFT 拍卖 API 测试脚本

BASE_URL="http://localhost:8080"

echo "======================================"
echo "NFT 拍卖 API 测试"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 健康检查
echo -e "${YELLOW}1. 健康检查${NC}"
curl -s "$BASE_URL/health" | jq
echo ""
echo ""

# 2. 获取拍卖列表
echo -e "${YELLOW}2. 获取拍卖列表（前10个）${NC}"
curl -s "$BASE_URL/api/auctions?page=1&page_size=10" | jq
echo ""
echo ""

# 3. 获取活跃拍卖列表
echo -e "${YELLOW}3. 获取活跃拍卖列表${NC}"
curl -s "$BASE_URL/api/auctions?status=active" | jq
echo ""
echo ""

# 4. 获取拍卖详情
echo -e "${YELLOW}4. 获取拍卖详情（ID=0）${NC}"
read -p "请输入拍卖ID [0]: " auction_id
auction_id=${auction_id:-0}
curl -s "$BASE_URL/api/auctions/$auction_id" | jq
echo ""
echo ""

# 5. 获取拍卖的出价历史
echo -e "${YELLOW}5. 获取拍卖的出价历史${NC}"
curl -s "$BASE_URL/api/auctions/$auction_id/bids?page=1&page_size=10" | jq
echo ""
echo ""

# 6. 从合约读取拍卖信息
echo -e "${YELLOW}6. 从合约读取拍卖信息${NC}"
curl -s "$BASE_URL/api/auctions/$auction_id/contract" | jq
echo ""
echo ""

# 7. 获取统计信息
echo -e "${YELLOW}7. 获取统计信息${NC}"
curl -s "$BASE_URL/api/stats" | jq
echo ""
echo ""

# 8. 参与出价（需要私钥）
echo -e "${YELLOW}8. 参与出价（需要私钥）${NC}"
read -p "是否要测试出价功能？(y/n) [n]: " test_bid
if [[ $test_bid == "y" ]]; then
    read -p "请输入拍卖ID: " bid_auction_id
    read -p "请输入出价金额（wei）[2000000000000000000]: " bid_amount
    bid_amount=${bid_amount:-2000000000000000000}
    read -p "请输入代币地址（空=ETH）: " token_address
    token_address=${token_address:-0x0000000000000000000000000000000000000000}
    read -sp "请输入私钥（不带0x）: " private_key
    echo ""
    
    curl -s -X POST "$BASE_URL/api/auctions/$bid_auction_id/bid" \
        -H "Content-Type: application/json" \
        -d "{
            \"amount\": \"$bid_amount\",
            \"token_address\": \"$token_address\",
            \"private_key\": \"$private_key\"
        }" | jq
    echo ""
fi
echo ""

# 9. 结束拍卖（需要私钥）
echo -e "${YELLOW}9. 结束拍卖（需要私钥）${NC}"
read -p "是否要测试结束拍卖功能？(y/n) [n]: " test_end
if [[ $test_end == "y" ]]; then
    read -p "请输入拍卖ID: " end_auction_id
    read -sp "请输入私钥（不带0x）: " private_key
    echo ""
    
    curl -s -X POST "$BASE_URL/api/auctions/$end_auction_id/end" \
        -H "Content-Type: application/json" \
        -d "{
            \"private_key\": \"$private_key\"
        }" | jq
    echo ""
fi
echo ""

# 10. 获取用户出价记录
echo -e "${YELLOW}10. 获取用户出价记录${NC}"
read -p "请输入用户地址 [跳过]: " bidder_address
if [[ ! -z "$bidder_address" ]]; then
    curl -s "$BASE_URL/api/bids?bidder=$bidder_address" | jq
    echo ""
fi

echo ""
echo -e "${GREEN}======================================"
echo "测试完成！"
echo "======================================${NC}"
