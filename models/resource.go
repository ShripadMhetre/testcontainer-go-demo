package models

type CreateResourceRequest struct {
	SiteGeoLocation string    `json:"site_geo_location,omitempty"`
	AssetInfo       AssetInfo `json:"asset_info" binding:"required"`
}

type AssetInfo struct {
	AssetTag    string `json:"asset_tag" binding:"required,min=1,max=16"`
	AssetType   string `json:"asset_type" binding:"omitempty,min=1,max=16"`
	AssetFamily string `json:"asset_family" binding:"omitempty,min=1,max=16"`
	ServerType  string `json:"server_type" binding:"omitempty,min=1,max=16"`
}

type Resource struct {
	ID              string    `gorm:"type:varchar;primaryKey;not null" json:"resource_id"`
	OfferId         string    `json:"offer_id,omitempty"`
	SiteGeoLocation string    `json:"site_geo_location,omitempty"`
	AssetInfo       AssetInfo `json:"asset_info" gorm:"embedded"`
}
