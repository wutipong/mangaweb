package meta

type SearchField string
type SortField string
type SortOrder string

const (
	SearchFieldName     = SearchField("name")
	SearchFieldFavorite = SearchField("favorite")
	SearchFieldTag      = SearchField("tag")

	SortFieldName       = SortField("name")
	SortFieldCreateTime = SortField("createTime")

	SortOrderAscending  = SortOrder("ascending")
	SortOrderDescending = SortOrder("descending")
)

type SearchCriteria struct {
	Field SearchField
	Value interface{}
}
