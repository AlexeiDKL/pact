package garant

type Documents struct {
	Documents []Document
}
type Document struct {
	URL   string `json:"url"`
	Topic int    `json:"topic"`
	Name  string `json:"name"`
}

type Data struct {
	Text      string   `json:"text"`
	Count     int      `json:"count"`
	IsQuery   bool     `json:"isQuery"`
	Kind      []string `json:"kind"`
	Sort      int      `json:"sort"`
	SortOrder int      `json:"sortOrder"`
}
