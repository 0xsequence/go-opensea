package opensea

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/0xsequence/go-sequence/lib/prototyp"
	"github.com/cespare/xxhash/v2"
	"github.com/goware/cachestore"
	"github.com/goware/cachestore/cachestorectl"
	"github.com/goware/cachestore/memlru"
	"github.com/goware/superr"
	"github.com/rs/zerolog"
)

var OpenSeaBaseURL = "https://api.opensea.io/api/v1"

const (
	DefaultCacheSize           = 10_000 // only used with memlru
	DefaultCacheExpiry         = 1 * time.Hour
	DefaultCacheNotFoundExpiry = 24 * time.Hour
	DefaultTimeout             = 10 * time.Second
)

type OpenSea struct {
	log     zerolog.Logger
	client  *http.Client
	apiKey  string
	cache   cachestore.Store[*apiResponse]
	timeout time.Duration
}

type Option struct {
	Timeout      time.Duration
	CacheBackend cachestore.Backend
	// CacheExpiry  *time.Duration // TODO ..
}

var (
	ErrTimeout      = errors.New("opensea: request timeout")
	ErrRateLimited  = errors.New("opensea: rate-limited by service")
	ErrUnauthorized = errors.New("opensea: unauthorized")
	ErrFail         = errors.New("opensea: fail")
)

func NewOpenSeaClient(log zerolog.Logger, client *http.Client, APIKey string, opts ...Option) (*OpenSea, error) {
	var cache cachestore.Store[*apiResponse]
	var err error

	// TODO: get options CacheExpiry for now, we'll just use DefaultCacheExpiry. But would be nice for Options
	// to have an override

	if len(opts) == 0 || opts[0].CacheBackend == nil {
		cache, err = memlru.NewWithSize[*apiResponse](DefaultCacheSize, cachestore.WithDefaultKeyExpiry(DefaultCacheExpiry))
	} else {
		cache, err = cachestorectl.Open[*apiResponse](opts[0].CacheBackend, cachestore.WithDefaultKeyExpiry(DefaultCacheExpiry))
	}
	if err != nil {
		return nil, fmt.Errorf("NewOpenSeaClient failed to open cache: %w", err)
	}

	if client == nil {
		client = http.DefaultClient
	}

	var timeout time.Duration
	for _, opt := range opts {
		if opt.Timeout > 0 {
			timeout = opt.Timeout
		}
	}
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	return &OpenSea{
		log:     log,
		client:  client,
		apiKey:  APIKey,
		cache:   cache,
		timeout: timeout,
	}, nil
}

func (o *OpenSea) RetrieveAsset(ctx context.Context, chainID uint64, tokenId prototyp.BigInt, assetContractAddress string) (*Asset, error) {
	if tokenId.Int().Sign() < 0 || assetContractAddress == "" {
		return nil, fmt.Errorf("tokenId and assetContractAddress are required")
	}

	url := fmt.Sprintf("%s/%s/%s", apiAssetEndpoint, assetContractAddress, tokenId.String())

	var assetsResponse *Asset
	_, err := o.doRequest(ctx, "GET", url, nil, &assetsResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve asset: %w", err)
	}

	return assetsResponse, nil
}

func (o *OpenSea) RetrieveAssetListing(ctx context.Context, chainID uint64, tokenId int64, assetContracAddress string) (*GetListingsResponse, error) {
	if tokenId < 0 || assetContracAddress == "" {
		return nil, fmt.Errorf("tokenId and assetContractAddress are required")
	}

	url := fmt.Sprintf("%s/%s/%d%s", apiAssetEndpoint, assetContracAddress, tokenId, apiListingEndpoint)

	var listingsResponse *GetListingsResponse
	_, err := o.doRequest(ctx, "GET", url, nil, &listingsResponse)
	if err != nil {
		o.log.Warn().Msgf("failed to retrieve listing: %v", err)
		return nil, err
	}

	return listingsResponse, nil
}

func (o *OpenSea) doRequest(ctx context.Context, httpMethod, endpointURL string, in, out interface{}) (int, error) {
	endpointPath := OpenSeaBaseURL + endpointURL
	response, present, err := o.cache.Get(ctx, apiRequestKey(&apiRequest{httpMethod: httpMethod, endpointPath: endpointPath}))
	if present && err == nil {
		err = json.Unmarshal(response.outBody, &out)
		if err != nil {
			o.log.Debug().Err(err).Msgf("Error unmarshalling response from cache")
		}

		// return successes
		if response.statusCode >= 200 && response.statusCode <= 299 {
			return response.statusCode, nil
		}

		// return cached failures
		if response.statusCode != 401 && response.statusCode != 429 {
			return response.statusCode, superr.Wrap(ErrFail, fmt.Errorf("opensea fail, status code %d: %s", response.statusCode, response.outBody))
		}
	}

	client := o.client
	retryCount := 0

	requestStarted := time.Now()

retry:
	var reqBody io.Reader

	// Below we can have a max amount of time, aka considering rate-limits too, after which
	// point we still stop and just return the request timed out. This is okay but not amazing.
	if time.Since(requestStarted) > o.timeout && o.timeout > 0 {
		return http.StatusRequestTimeout, superr.Wrap(ErrTimeout, ErrRateLimited)
	}

	o.log.Info().Msgf("opensea doRequest to endpoint %s", endpointURL)

	if in != nil {
		reqBodyBytes, err := json.Marshal(in)
		if err != nil {
			return 0, fmt.Errorf("failed to marshal json request: %w", err)
		}
		reqBody = bytes.NewBuffer(reqBodyBytes)
	}

	req, err := http.NewRequestWithContext(ctx, httpMethod, endpointPath, reqBody)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", o.apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return 0, superr.Wrap(ErrFail, err)
	}
	defer resp.Body.Close()

	if err = ctx.Err(); err != nil {
		return 0, fmt.Errorf("aborted because context was done: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		if resp.StatusCode == 429 && retryCount < 20 {
			select {
			case <-ctx.Done():
				// done
			default:
				retryCount++
				delay := time.Duration(retryCount) * time.Second * 2
				o.log.Warn().Msgf("opensea request to endpoint %s hit rate-limit, delaying for %s", endpointURL, delay)
				time.Sleep(delay)
				goto retry
			}
		}

		if resp.StatusCode == http.StatusUnauthorized {
			return resp.StatusCode, superr.Wrap(ErrUnauthorized, fmt.Errorf("invalid or expired API key"))
		}

		if resp.StatusCode == http.StatusNotFound {
			err = o.cache.SetEx(ctx, apiRequestKey(&apiRequest{httpMethod: httpMethod, endpointPath: endpointPath}), &apiResponse{
				statusCode: resp.StatusCode,
			}, DefaultCacheNotFoundExpiry)
			if err != nil {
				o.log.Debug().Err(err).Msgf("Error setting cache")
			}
		}

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, superr.Wrap(ErrFail, fmt.Errorf("failed to read response body: %w", err))
		}
		err = o.cache.Set(ctx, apiRequestKey(&apiRequest{httpMethod: httpMethod, endpointPath: endpointPath}), &apiResponse{
			statusCode: resp.StatusCode,
			outBody:    respBody,
		})
		if err != nil {
			o.log.Debug().Err(err).Msgf("Error setting cache")
		}

		return resp.StatusCode, superr.Wrap(ErrFail, fmt.Errorf("status code %d: %s", resp.StatusCode, resp.Body))
	}

	var respBody []byte

	if out != nil {
		respBody, err = io.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, fmt.Errorf("failed to read response body: %w", err)
		}

		err = json.Unmarshal(respBody, &out)
		if err != nil {
			return resp.StatusCode, fmt.Errorf("failed to unmarshal json response body: %w", err)
		}
		if err = ctx.Err(); err != nil {
			return resp.StatusCode, fmt.Errorf("aborted because context was done: %w", err)
		}
	}

	err = o.cache.Set(ctx, apiRequestKey(&apiRequest{httpMethod: httpMethod, endpointPath: endpointPath}), &apiResponse{
		statusCode: resp.StatusCode,
		outBody:    respBody,
	})
	if err != nil {
		o.log.Debug().Err(err).Msgf("Error setting cache")
	}

	return resp.StatusCode, nil
}

func (o *OpenSea) SetHttpClient(httpClient *http.Client) {
	o.client = httpClient
}

type apiRequest struct {
	httpMethod      string
	endpointPath    string
	queryParameters url.Values
	inBody          []byte
}

type apiResponse struct {
	statusCode int
	outBody    []byte
	err        error
}

func apiRequestKey(req *apiRequest) string {
	return fmt.Sprintf("opensea-req:%d", xxhash.Sum64(append([]byte(req.httpMethod+req.endpointPath+req.queryParameters.Encode()), req.inBody...)))
}
