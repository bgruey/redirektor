package repo

import (
	"redirektor/server/model"

	"errors"
	"gorm.io/gorm"
)

func (pc *PostgresClient) CreateApiKey(key *model.ApiKey, tx *gorm.DB) error {
	tx = pc.selfIfNil(tx)

	tx.Create(key)
	return tx.Error
}

func (pc *PostgresClient) GetApiKey(key string, tx *gorm.DB) (retKey *model.ApiKey, err error) {
	tx = pc.selfIfNil(tx)

	retKey = &model.ApiKey{}

	err = tx.Where("key = ? and (deleted_at is null or deleted_at > extract(epoch from now()))", key).First(&retKey).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (pc *PostgresClient) CountApiKeys(tx *gorm.DB) (count int64, err error) {
	tx = pc.selfIfNil(tx)

	err = tx.Model(model.ApiKey{}).Count(&count).Error

	return
}

func (pc *PostgresClient) GetSingleApiKey(tx *gorm.DB) (key *model.ApiKey, err error) {
	tx = pc.selfIfNil(tx)
	key = &model.ApiKey{}

	err = tx.Model(key).First(&key).Error
	return
}
