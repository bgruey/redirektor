package repo

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"

	"redirektor/server/model"
	"redirektor/server/utils"
)

func (pc *PostgresClient) CreateRedirect(redirect *model.Redirect, tx *gorm.DB) (err error) {
	tx = pc.selfIfNil(tx)

	linkRedirect, err := pc.GetRedirectByLink(redirect.Link, tx)
	if err != nil {
		return err
	}

	if linkRedirect != nil {
		log.Printf("Found link: %s --> %s", linkRedirect.Link, linkRedirect.Hash)
		// for response
		redirect.Hash = linkRedirect.Hash
		redirect.QRCode = linkRedirect.QRCode
		// trigger update instead of create on save
		redirect.ID = linkRedirect.ID
		return
	}

	hash := utils.Sha256Base64(redirect.Link)

	chars := 1
	for {
		redirect.Hash = hash[0:chars]
		existingLink, err := pc.GetRedirectByHash(redirect.Hash, tx, false)
		if err != nil {
			return err
		}
		if existingLink == nil {
			break
		}
		chars++
	}

	err = tx.Create(redirect).Error

	return
}

func (pc *PostgresClient) GetRedirectByHash(hash string, tx *gorm.DB, lock bool) (redirect *model.Redirect, err error) {
	tx = pc.selfIfNil(tx)

	redirect = &model.Redirect{}

	if lock {
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("hash = ?", hash).First(redirect).Error
	} else {
		err = tx.Where("hash = ?", hash).First(redirect).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return redirect, nil
}

func (pc *PostgresClient) GetRedirectByLink(link string, tx *gorm.DB) (redirect *model.Redirect, err error) {
	tx = pc.selfIfNil(tx)

	redirect = &model.Redirect{}

	err = tx.Where("link = ?", link).First(redirect).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return redirect, nil
}

func (pc *PostgresClient) GetIncrementRedirectByHash(hash string) (link string, err error) {
	tx := pc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}

		tx.Commit()
	}()

	redirect, err := pc.GetRedirectByHash(hash, tx, true)
	if err != nil {
		panic(err)
	}
	if redirect == nil {
		return "", nil
	}

	redirect.Count++
	err = tx.Save(redirect).Error

	return redirect.Link, err
}

func (pc *PostgresClient) SaveRedirect(redirect *model.Redirect, tx *gorm.DB) (err error) {
	tx = pc.selfIfNil(tx)

	err = tx.Save(redirect).Error

	return
}
