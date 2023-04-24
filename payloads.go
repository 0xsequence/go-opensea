package opensea

type GetListingsResponse struct {
	SeaPortListings []*Listings `json:"seaport_listings"`
}

type CollectionStatsResponse struct {
	Stats *CollectionStats `json:"stats"`
}
