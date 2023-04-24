package opensea

//
// Account
//

type User struct {
	Username string `json:"username"`
}

type Address string

type Account struct {
	User          User    `json:"user"`
	ProfileImgURL string  `json:"profile_img_url"`
	Address       Address `json:"address"`
	Config        string  `json:"config"`
	DiscordID     string  `json:"discord_id"`
}

//
// Asset
//

type Asset struct {
	ID                   int64          `json:"id"`
	TokenID              string         `json:"token_id"`
	NumSales             int64          `json:"num_sales"`
	BackgroundColor      string         `json:"background_color"`
	ImageURL             string         `json:"image_url"`
	ImagePreviewURL      string         `json:"image_preview_url"`
	ImageThumbnailURL    string         `json:"image_thumbnail_url"`
	ImageOriginalURL     string         `json:"image_original_url"`
	AnimationURL         string         `json:"animation_url"`
	AnimationOriginalURL string         `json:"animation_original_url"`
	Name                 string         `json:"name"`
	Description          string         `json:"description"`
	ExternalLink         string         `json:"external_link"`
	AssetContract        *AssetContract `json:"asset_contract"`
	Owner                *Account       `json:"owner"`
	Permalink            string         `json:"permalink"`
	Collection           *Collection    `json:"collection"`
	Decimals             int64          `json:"decimals"`
	TokenMetadata        string         `json:"token_metadata"`
	Traits               interface{}    `json:"traits"`
	LastSale             LastSale       `json:"last_sale"`
}

type AssetContract struct {
	Address                     string `json:"address"`
	AssetContractType           string `json:"asset_contract_type"`
	CreatedDate                 string `json:"created_date"`
	Name                        string `json:"name"`
	NftVersion                  string `json:"nft_version"`
	OpenseaVersion              string `json:"opensea_version"`
	Owner                       int    `json:"owner"`
	SchemaName                  string `json:"schema_name"`
	Symbol                      string `json:"symbol"`
	TotalSupply                 string `json:"total_supply"`
	Description                 string `json:"description"`
	ExternalLink                string `json:"external_link"`
	ImageURL                    string `json:"image_url"`
	DefaultToFiat               bool   `json:"default_to_fiat"`
	DevBuyerFeeBasisPoints      int    `json:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints     int    `json:"dev_seller_fee_basis_points"`
	OnlyProxiedTransfers        bool   `json:"only_proxied_transfers"`
	OpenseaBuyerFeeBasisPoints  int    `json:"opensea_buyer_fee_basis_points"`
	OpenseaSellerFeeBasisPoints int    `json:"opensea_seller_fee_basis_points"`
	BuyerFeeBasisPoints         int    `json:"buyer_fee_basis_points"`
	SellerFeeBasisPoints        int    `json:"seller_fee_basis_points"`
	PayoutAddress               string `json:"payout_address"`
}

type LastSale struct {
	TotalPrice   string       `json:"total_price"`
	PaymentToken PaymentToken `json:"payment_token"`
	Quantity     string       `json:"quantity"`
}

type PaymentToken struct {
	ID       int    `json:"id"`
	Symbol   string `json:"symbol"`
	Address  string `json:"address"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	EthPrice string `json:"eth_price"`
	UsdPrice string `json:"usd_price"`
}

//
// Collection
//

type Collection struct {
	BannerImageUrl              string      `json:"banner_image_url" bson:"banner_image_url"`
	ChatUrl                     string      `json:"chat_url" bson:"chat_url"`
	CreatedDate                 string      `json:"created_date" bson:"created_date"`
	DefaultToFiat               bool        `json:"default_to_fiat" bson:"default_to_fiat"`
	Description                 string      `json:"description" bson:"description"`
	DevBuyerFeeBasisPoints      string      `json:"dev_buyer_fee_basis_points" bson:"dev_buyer_fee_basis_points"`
	DevSellerFeeBasisPoints     string      `json:"dev_seller_fee_basis_points" bson:"dev_seller_fee_basis_points"`
	DiscordUrl                  string      `json:"discord_url" bson:"discord_url"`
	DisplayData                 interface{} `json:"display_data" bson:"display_data"`
	ExternalUrl                 string      `json:"external_url" bson:"external_url"`
	Featured                    bool        `json:"featured" bson:"featured"`
	FeaturedImageUrl            string      `json:"featured_image_url" bson:"featured_image_url"`
	Hidden                      bool        `json:"hidden" bson:"hidden"`
	SafeListRequestStatus       string      `json:"safelist_request_status" bson:"safelist_request_status"`
	ImageUrl                    string      `json:"image_url" bson:"image_url"`
	IsSubjectToWhitelist        bool        `json:"is_subject_to_whitelist" bson:"is_subject_to_whitelist"`
	LargeImageUrl               string      `json:"large_image_url" bson:"large_image_url"`
	MediumUsername              string      `json:"medium_username" bson:"medium_username"`
	Name                        string      `json:"name" bson:"name"`
	OnlyProxiedTransfers        bool        `json:"only_proxied_transfers" bson:"only_proxied_transfers"`
	OpenseaBuyerFeeBasisPoints  string      `json:"opensea_buyer_fee_basis_points" bson:"opensea_buyer_fee_basis_points"`
	OpenseaSellerFeeBasisPoints string      `json:"opensea_seller_fee_basis_points" bson:"opensea_seller_fee_basis_points"`
	PayoutAddress               string      `json:"payout_address" bson:"payout_address"`
	RequireEmail                bool        `json:"require_email" bson:"require_email"`
	ShortDescription            string      `json:"short_description" bson:"short_description"`
	Slug                        string      `json:"slug" bson:"slug"`
	TelegramUrl                 string      `json:"telegram_url" bson:"telegram_url"`
	TwitterUsername             string      `json:"twitter_username" bson:"twitter_username"`
	InstagramUsername           string      `json:"instagram_username" bson:"instagram_username"`
	WikiUrl                     string      `json:"wiki_url" bson:"wiki_url"`
}

type CollectionStats struct {
	OneDayVolume          float64 `json:"one_day_volume"`
	OneDayChange          float64 `json:"one_day_change"`
	OneDaySales           float64 `json:"one_day_sales"`
	OneDayAveragePrice    float64 `json:"one_day_average_price"`
	SevenDayVolume        float64 `json:"seven_day_volume"`
	SevenDayChange        float64 `json:"seven_day_change"`
	SevenDaySales         float64 `json:"seven_day_sales"`
	SevenDayAveragePrice  float64 `json:"seven_day_average_price"`
	ThirtyDayVolume       float64 `json:"thirty_day_volume"`
	ThirtyDayChange       float64 `json:"thirty_day_change"`
	ThirtyDaySales        float64 `json:"thirty_day_sales"`
	ThirtyDayAveragePrice float64 `json:"thirty_day_average_price"`
	TotalVolume           float64 `json:"total_volume"`
	TotalSales            float64 `json:"total_sales"`
	TotalSupply           float64 `json:"total_supply"`
	Count                 float64 `json:"count"`
	NumOwners             float64 `json:"num_owners"`
	AveragePrice          float64 `json:"average_price"`
	NumReports            float64 `json:"num_reports"`
	MarketCap             float64 `json:"market_cap"`
	FloorPrice            float64 `json:"floor_price"`
}

//
// Listing
//

type Listings struct {
	CurrentPrice string `json:"current_price"`
}

type PaymentTokenContract struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Decimals int    `json:"decimals"`
	EthPrice string `json:"eth_price"`
	UsdPrice string `json:"usd_price"`
}
