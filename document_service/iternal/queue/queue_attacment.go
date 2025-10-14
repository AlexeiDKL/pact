package queue

type QueueAttachmentSend struct {
	Text      string   `json:"text"`
	Count     int      `json:"count"`
	IsQuery   bool     `json:"isQuery"`
	Kind      []string `json:"kind"`
	Sort      int      `json:"sort"`
	SortOrder int      `json:"sortOrder"`
}

type QueueAttachmentRequest struct {
	Documents []Document `json:"documents"`
}

type Document struct {
	Url   string `json:"url"`
	Topic int    `json:"topic"`
	Name  string `json:"name"`
}

func (q *QueueAttachmentSend) CreateQuery(name string) QueueAttachmentSend {
	return QueueAttachmentSend{
		Text:      "& BOOL(& MorphoName(" + name + "))",
		Count:     10,
		IsQuery:   true,
		Kind:      []string{"001", "002", "003"},
		Sort:      0,
		SortOrder: 0,
	}
}
