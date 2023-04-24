package opensea

import (
	"context"
	"fmt"
)

func (o *OpenSea) RetrieveCollectionStats(ctx context.Context, chainID uint64, openSeaCollectionIdentifier string) (*CollectionStatsResponse, error) {
	if openSeaCollectionIdentifier == "" {
		return nil, fmt.Errorf("openSeaCollectionIdentifier is required")
	}

	url := fmt.Sprintf("%s/%s/%s", apiCollectionEndpoint, openSeaCollectionIdentifier, apiStatEndpoint)

	var collectionStatsResponse *CollectionStatsResponse
	_, err := o.doRequest(ctx, "GET", url, nil, &collectionStatsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve asset: %w", err)
	}

	return collectionStatsResponse, nil
}

func (o *OpenSea) RetrieveAssetContractInfo(ctx context.Context, chainID uint64, assetContractAddress string) (*AssetContract, error) {
	if assetContractAddress == "" {
		return nil, fmt.Errorf("assetContractAddress is required")
	}

	url := fmt.Sprintf("%s/%s", apiAssetContractEndpoint, assetContractAddress)

	var assetContractResponse *AssetContract
	_, err := o.doRequest(ctx, "GET", url, nil, &assetContractResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve asset: %w", err)
	}

	return assetContractResponse, nil
}
