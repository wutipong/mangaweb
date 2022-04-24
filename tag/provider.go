package tag

type Provider interface {
	Close() error
	Delete(t Tag) error
	IsTagExist(name string) bool
	Read(name string) (t Tag, err error)
	ReadAll() (tags []Tag, err error)
	Write(t Tag) error
}

type ProviderFactory func() (p Provider, err error)
