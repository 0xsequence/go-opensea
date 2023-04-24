package opensea

import (
	"context"
	"net/http"
	"testing"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/rs/zerolog"
	"github.com/seborama/govcr"
	"github.com/stretchr/testify/assert"
)

var (
	openSeaClient *OpenSea
	client        *http.Client
	err           error
)

const APIKey = "08826c83664f4c5585bc5b452abe7d9e"

func init() {
	client = &http.Client{}
	openSeaClient, err = NewOpenSeaClient(zerolog.Logger{}, client, APIKey)
	if err != nil {
		panic(err)
	}
}

func TestGetAssets(t *testing.T) {
	vcr := govcr.NewVCR("openseaGetAssets", &govcr.VCRConfig{
		LongPlay:  true,
		RemoveTLS: true,
	})
	openSeaClient.SetHttpClient(vcr.Client)

	ctx := context.Background()
	assets, err := openSeaClient.RetrieveAsset(ctx, 1, prototyp.NewBigInt(9257), "0x282bdd42f4eb70e7a9d9f40c8fea0825b7f68c5d")
	assert.NoError(t, err)
	assert.NotNil(t, assets)
	assert.IsType(t, &Asset{}, assets)
}

func TestGetListings(t *testing.T) {
	vcr := govcr.NewVCR("openseaGetListings", &govcr.VCRConfig{
		LongPlay:  true,
		RemoveTLS: true,
	})
	openSeaClient.SetHttpClient(vcr.Client)

	ctx := context.Background()
	listings, err := openSeaClient.RetrieveAssetListing(ctx, 1, 2809, "0x2d0d57d004f82e9f4471caa8b9f8b1965a814154")
	assert.NoError(t, err)
	assert.NotNil(t, listings)
	assert.IsType(t, &GetListingsResponse{}, listings)
}

func TestCollectionStats(t *testing.T) {
	ctx := context.Background()
	stats, err := openSeaClient.RetrieveCollectionStats(ctx, 1, "doodles-official")
	assert.NoError(t, err)
	assert.NotNil(t, stats)
	assert.IsType(t, &CollectionStatsResponse{}, stats)
}

func TestAssetContract(t *testing.T) {
	ctx := context.Background()
	contractInfo, err := openSeaClient.RetrieveAssetContractInfo(ctx, 137, "0x631998e91476DA5B870D741192fc5Cbc55F5a52E")
	assert.NoError(t, err)
	assert.NotNil(t, contractInfo)
	assert.IsType(t, &AssetContract{}, contractInfo)
}
