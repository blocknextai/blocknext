package cache

type cachedContentResponse struct {
	Name       string `json:"name"`
	ExpireTime string `json:"expireTime"`
}
