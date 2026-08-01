package getnotificationcounts

type GetNotificationCountsResponse struct {
	Unread int64 `json:"unread"`
	Unseen int64 `json:"unseen"`
}
