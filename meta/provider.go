package meta

//Provider meta data provider.
type Provider interface {
	IsItemExist(name string) bool
	Write(i Item) error
	New(name string) (i Item, err error)
	Delete(i Item) error
	Read(name string) (i Item, err error)
	Open(name string) (i Item, err error)
	ReadAll() (items []Item, err error)
	Close() error
}