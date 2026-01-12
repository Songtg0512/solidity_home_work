package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OpenSeaService OpenSea API 服务
type OpenSeaService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewOpenSeaService 创建 OpenSea 服务实例
func NewOpenSeaService(apiKey string) *OpenSeaService {
	return &OpenSeaService{
		apiKey:  apiKey,
		baseURL: "https://api.opensea.io/api/v2",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CollectionStats 集合统计信息
type CollectionStats struct {
	FloorPrice float64 `json:"floor_price"`
	TotalVolume float64 `json:"total_volume"`
	OneDayVolume float64 `json:"one_day_volume"`
	SevenDayVolume float64 `json:"seven_day_volume"`
	ThirtyDayVolume float64 `json:"thirty_day_volume"`
	TotalSales int `json:"total_sales"`
	Count int `json:"count"`
}

// OpenSeaCollection OpenSea 集合信息
type OpenSeaCollection struct {
	Collection string `json:"collection"`
	Name string `json:"name"`
	Description string `json:"description"`
	ImageURL string `json:"image_url"`
	BannerImageURL string `json:"banner_image_url"`
	Stats CollectionStats `json:"stats"`
}

// GetCollectionStats 获取集合统计信息（包括地板价）
func (s *OpenSeaService) GetCollectionStats(collectionSlug string) (*CollectionStats, error) {
	url := fmt.Sprintf("%s/collections/%s/stats", s.baseURL, collectionSlug)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if s.apiKey != "" {
		req.Header.Set("X-API-KEY", s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("opensea API error: %s, body: %s", resp.Status, string(body))
	}

	var result struct {
		Stats CollectionStats `json:"stats"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result.Stats, nil
}

// GetCollection 获取集合信息
func (s *OpenSeaService) GetCollection(collectionSlug string) (*OpenSeaCollection, error) {
	url := fmt.Sprintf("%s/collections/%s", s.baseURL, collectionSlug)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if s.apiKey != "" {
		req.Header.Set("X-API-KEY", s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("opensea API error: %s, body: %s", resp.Status, string(body))
	}

	var result OpenSeaCollection
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetCollectionByContract 通过合约地址获取集合 slug
func (s *OpenSeaService) GetCollectionByContract(contractAddress string) (string, error) {
	// OpenSea v2 API 需要通过合约地址查询集合
	url := fmt.Sprintf("%s/chain/ethereum/contract/%s", s.baseURL, strings.ToLower(contractAddress))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	if s.apiKey != "" {
		req.Header.Set("X-API-KEY", s.apiKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("opensea API error: %s, body: %s", resp.Status, string(body))
	}

	var result struct {
		Collection string `json:"collection"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Collection, nil
}

// GetFloorPriceByContract 通过合约地址获取地板价
func (s *OpenSeaService) GetFloorPriceByContract(contractAddress string) (float64, error) {
	// 先获取集合 slug
	collectionSlug, err := s.GetCollectionByContract(contractAddress)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection slug: %w", err)
	}

	// 再获取统计信息
	stats, err := s.GetCollectionStats(collectionSlug)
	if err != nil {
		return 0, fmt.Errorf("failed to get collection stats: %w", err)
	}

	return stats.FloorPrice, nil
}
