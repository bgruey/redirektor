package model

type Redirect struct {
	ID        int64  `gorm:"unique;primaryKey;autoIncrement" json:"-"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"-"`
	Hash      string `gorm:"unique_index"`
	Link      string `gorm:"unique_index" json:"link"`
	Count     int64  `gorm:"default:0"`
}

func NewLink(link string) *Redirect {
	ret := new(Redirect)
	ret.Link = link

	return ret
}
