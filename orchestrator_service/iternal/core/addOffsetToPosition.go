package core

type Item struct {
	Name     string `json:"name"`
	Caption  string `json:"caption"`
	Start    int    `json:"startPosition"`
	End      int    `json:"endPosition"`
	Children Child  `json:"children"`
}

type Child struct {
	Item []Item `json:"Item"`
}

type Root struct {
	Item []Item `json:"Item"`
}

func AddOffsetToPositions(items []Item, offset int) {
	for i := range items {
		items[i].Start += offset
		items[i].End += offset
		// Рекурсивно обрабатываем вложенные элементы
		AddOffsetToPositions(items[i].Children.Item, offset)
	}
}
