package response

type Namespace struct {
	Name              string `json:"name"`
	CreationTimestamp int64  `json:"create_timestamp"`
	Status            string `json:"status"`
}
