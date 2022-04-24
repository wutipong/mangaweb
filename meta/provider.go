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

//Provider meta data provider.
type Provider interface {
	IsItemExist(name string) bool
	Write(i Meta) error
	Delete(i Meta) error
	Read(name string) (i Meta, err error)
	Open(name string) (i Meta, err error)
	ReadAll() (items []Meta, err error)
	Search(criteria []SearchCriteria, sort SortField, order SortOrder, pageSize int, page int) (items []Meta, err error)
	Count(criteria []SearchCriteria) (count int64, err error)
	Close() error
}
