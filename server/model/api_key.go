package model

import (
	"redirektor/server/utils"

	"github.com/google/uuid"
)

type ApiKey struct {
	ID        int64  `gorm:"unique;primaryKey;autoIncrement" json:"-"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt int64  `gorm:"autoUpdateTime:milli" json:"-"`
	DeletedAt int64  `gorm:"default: null"`
	Key       string `gorm:"unique_index"`
}

func NewApiKey() *ApiKey {
	key := uuid.New()
	ret := new(ApiKey)
	ret.Key = utils.Sha256Base64(key.String())

	return ret
}
