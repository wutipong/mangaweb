package meta

type SearchField string
type SortField string
type SortOrder string

const (
	SearchFieldName     = SearchField("name")
	SearchFieldFavorite = SearchField("favorite")

	SortFieldName       = SortField("name")
	SortFieldCreateTime = SortField("createTime")

	SortOrderAscending  = SortOrder("ascending")
	SortOrderDescending = SortOrder("descending")
)

type SearchCriteria struct {
	Field SearchField
	Value interface{}
}

//Provider meta data provider.
type Provider interface {
	IsItemExist(name string) bool
	Write(i Item) error
	Delete(i Item) error
	Read(name string) (i Item, err error)
	Open(name string) (i Item, err error)
	ReadAll() (items []Item, err error)
	Search(criteria []SearchCriteria, sort SortField, order SortOrder, pageSize int, page int) (items []Item, err error)
	Count(criteria []SearchCriteria) (count int64, err error)
	Close() error
}
