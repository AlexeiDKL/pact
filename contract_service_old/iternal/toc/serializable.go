package toc

type TOCItemWithChildren struct {
	Name         string                `json:"name"`
	Caption      string                `json:"caption"`
	StartPos     int                   `json:"startPosition"`
	EndPos       int                   `json:"endPosition"`
	ChildrenWrap *ChildrenArrayWrapper `json:"children,omitempty"`
}

type ChildrenArrayWrapper struct {
	Item []TOCItemWithChildren `json:"Item"`
}

// Рекурсивное преобразование
func ConvertTOCItem(item TOCItem) TOCItemWithChildren {
	var children []TOCItemWithChildren
	for _, child := range item.Children {
		children = append(children, ConvertTOCItem(child))
	}
	var wrap *ChildrenArrayWrapper
	if len(children) > 0 {
		wrap = &ChildrenArrayWrapper{Item: children}
	}
	return TOCItemWithChildren{
		Name:         item.Name,
		Caption:      item.Caption,
		StartPos:     item.StartPos,
		EndPos:       item.EndPos,
		ChildrenWrap: wrap,
	}
}

// Для корня:
// result := []TOCItemWithChildren{ConvertTOCItem(root1), ConvertTOCItem(root2)}
// json.Marshal(struct{Item []TOCItemWithChildren `json:"Item"`}{Item: result})
