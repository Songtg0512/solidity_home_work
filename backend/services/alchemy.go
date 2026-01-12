package services

import (
	"auction-backend/config"
	"auction-backend/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// AlchemyService Alchemy API 服务
type AlchemyService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewAlchemyService 创建 Alchemy 服务实例
func NewAlchemyService() *AlchemyService {
	return &AlchemyService{
		apiKey:  config.AppConfig.AlchemyAPIKey,
		baseURL: config.AppConfig.AlchemyBaseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// AlchemyNFT Alchemy NFT 响应
type AlchemyNFT struct {
	Contract struct {
		Address string `json:"address"`
	} `json:"contract"`
	TokenID  string `json:"id"`
	Title    string `json:"title"`
	Description string `json:"description"`
	TokenURI struct {
		Gateway string `json:"gateway"`
	} `json:"tokenUri"`
	Media []struct {
		Gateway string `json:"gateway"`
	} `json:"media"`
	Metadata struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Image       string `json:"image"`
		Attributes  []map[string]interface{} `json:"attributes"`
	} `json:"metadata"`
}

// AlchemyNFTsResponse Alchemy 获取 NFTs 响应
type AlchemyNFTsResponse struct {
	OwnedNfts []AlchemyNFT `json:"ownedNfts"`
	TotalCount int `json:"totalCount"`
	PageKey string `json:"pageKey,omitempty"`
}

// GetNFTsByOwner 获取某个地址拥有的所有 NFT
func (s *AlchemyService) GetNFTsByOwner(owner string, pageKey string) (*AlchemyNFTsResponse, error) {
	url := fmt.Sprintf("%s/getNFTs?owner=%s&withMetadata=true", s.baseURL, owner)
	if pageKey != "" {
		url += fmt.Sprintf("&pageKey=%s", pageKey)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("alchemy API error: %s, body: %s", resp.Status, string(body))
	}

	var result AlchemyNFTsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetNFTMetadata 获取 NFT 元数据
func (s *AlchemyService) GetNFTMetadata(contractAddress, tokenID string) (*models.NFTMetadata, error) {
	url := fmt.Sprintf("%s/getNFTMetadata?contractAddress=%s&tokenId=%s", s.baseURL, contractAddress, tokenID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("alchemy API error: %s, body: %s", resp.Status, string(body))
	}

	var alchemyNFT AlchemyNFT
	if err := json.NewDecoder(resp.Body).Decode(&alchemyNFT); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// 转换为 NFTMetadata 模型
	metadata := &models.NFTMetadata{
		Contract: strings.ToLower(contractAddress),
		TokenID:  tokenID,
		Name:     alchemyNFT.Metadata.Name,
		Description: alchemyNFT.Metadata.Description,
		Image:    alchemyNFT.Metadata.Image,
		LastSync: time.Now(),
	}

	// 序列化 attributes
	if len(alchemyNFT.Metadata.Attributes) > 0 {
		attrJSON, _ := json.Marshal(alchemyNFT.Metadata.Attributes)
		metadata.Attributes = string(attrJSON)
	}

	// 优先使用 media gateway
	if len(alchemyNFT.Media) > 0 && alchemyNFT.Media[0].Gateway != "" {
		metadata.Image = alchemyNFT.Media[0].Gateway
	}

	return metadata, nil
}

// GetContractMetadata 获取合约元数据
func (s *AlchemyService) GetContractMetadata(contractAddress string) (*models.NFTCollection, error) {
	url := fmt.Sprintf("%s/getContractMetadata?contractAddress=%s", s.baseURL, contractAddress)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("alchemy API error: %s, body: %s", resp.Status, string(body))
	}

	var result struct {
		Name        string `json:"name"`
		Symbol      string `json:"symbol"`
		TotalSupply string `json:"totalSupply"`
		TokenType   string `json:"tokenType"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	collection := &models.NFTCollection{
		Contract: strings.ToLower(contractAddress),
		Name:     result.Name,
		Symbol:   result.Symbol,
		LastSync: time.Now(),
	}

	return collection, nil
}
