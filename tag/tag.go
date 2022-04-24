package tag

const (
	CurrentVersion = 0
)

type Tag struct {
	Name      string `json:"name" db:"name" bson:"name"`
	Favorite  bool   `json:"favorite" db:"favorite" bson:"favorite"`
	Hidden    bool   `json:"hidden" db:"hidden" bson:"hidden"`
	Version   int    `json:"version" db:"version" bson:"version"`
	Thumbnail []byte `json:"thumbnail" db:"thumbnail" bson:"thumbnail"`
}

func NewTag(name string) Tag {
	return Tag{
		Name:     name,
		Favorite: false,
		Hidden:   false,
		Version:  CurrentVersion,
	}
}
