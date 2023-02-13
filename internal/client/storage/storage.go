package storage

import (
	"errors"
	"github.com/asdine/storm/v3"
	"github.com/rainset/gophkeeper/internal/client/model"
	"github.com/rainset/gophkeeper/pkg/logger"
)

type Base struct {
	user string
	db   *storm.DB
}

func New(databaseName string) (base *Base, err error) {

	db, err := storm.Open(databaseName)

	if err != nil {
		return base, err
	}

	base = &Base{db: db}

	return base, err
}

func (b *Base) SetUser(user string) {
	b.user = user
}

func (b *Base) Close() (err error) {
	err = b.db.Close()
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (b *Base) SetConfig(config model.ConfigStorage) (err error) {
	err = b.db.From(b.user).Set("store", "config", config)
	return err
}
func (b *Base) GetConfig() (c model.ConfigStorage, err error) {
	err = b.db.From(b.user).Get("store", "config", &c)
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return c, nil
		}
	}
	return c, err
}

func (b *Base) AddCard(card model.DataCard) (err error) {
	err = b.db.From(b.user).Save(&card)
	if err != nil {
		logger.Error()
	}
	return err
}
func (b *Base) GetCard(ID int) (card model.DataCard, err error) {
	err = b.db.From(b.user).One("LocalID", ID, &card)
	if err != nil {
		logger.Error()
	}
	return card, err
}
func (b *Base) GetAllCards() (cards []model.DataCard, err error) {
	err = b.db.From(b.user).All(&cards, storm.Reverse())
	if err != nil {
		logger.Error(err)
	}
	return cards, err
}
func (b *Base) DeleteCard(ID int) (err error) {

	var card model.DataCard
	card.LocalID = ID
	err = b.db.From(b.user).DeleteStruct(&card)
	if err != nil {
		logger.Error(err)
	}
	return err
}

func (b *Base) AddCred(cred model.DataCred) (err error) {
	err = b.db.From(b.user).Save(&cred)
	if err != nil {
		logger.Error()
	}
	return err
}
func (b *Base) GetCred(ID int) (cred model.DataCred, err error) {
	err = b.db.From(b.user).One("LocalID", ID, &cred)
	return cred, err
}
func (b *Base) GetAllCreds() (creds []model.DataCred, err error) {
	err = b.db.From(b.user).All(&creds, storm.Reverse())
	return creds, err
}
func (b *Base) DeleteCred(ID int) (err error) {
	var cred model.DataCred
	cred.LocalID = ID
	err = b.db.From(b.user).DeleteStruct(&cred)
	return err
}

func (b *Base) AddText(text model.DataText) (err error) {
	err = b.db.From(b.user).Save(&text)
	return err
}
func (b *Base) GetText(ID int) (text model.DataText, err error) {
	err = b.db.From(b.user).One("LocalID", ID, &text)
	return text, err
}
func (b *Base) GetAllTexts() (texts []model.DataText, err error) {
	err = b.db.From(b.user).All(&texts, storm.Reverse())
	return texts, err
}
func (b *Base) DeleteText(ID int) (err error) {

	var text model.DataText
	text.LocalID = ID
	err = b.db.From(b.user).DeleteStruct(&text)
	return err
}

func (b *Base) AddFile(file model.DataFile) (err error) {
	err = b.db.From(b.user).Save(&file)
	return err
}
func (b *Base) GetFile(ID int) (file model.DataFile, err error) {
	err = b.db.From(b.user).One("LocalID", ID, &file)
	return file, err
}
func (b *Base) GetAllFiles() (files []model.DataFile, err error) {
	err = b.db.From(b.user).All(&files, storm.Reverse())
	if err != nil {
		logger.Error(err)
	}
	return files, err
}
func (b *Base) DeleteFile(ID int) (err error) {

	var f model.DataFile
	f.LocalID = ID
	err = b.db.From(b.user).DeleteStruct(&f)
	if err != nil {
		logger.Error(err)
	}
	return err
}

//func (b *Base) CreateUser() (err error) {
//
//	var card model.DataCard
//	card, err = b.GetCard(1)
//	if err != nil {
//		logger.Error(err)
//	}
//
//	logger.Info(card.Number)
//
//	var cards []model.DataCard
//	err = b.db.All(&cards)
//	if err != nil {
//		logger.Error(err)
//	}
//
//	for i, v := range cards {
//		logger.Info("index: ", i, " ID ", v.ID, " num: ", v.Number)
//	}
//
//	//
//	//user := model.DataCard{
//	//	ExternalID: 11,
//	//	Title:      "Тинькоф тайтл",
//	//	Number:     "12312 123 123 213",
//	//	Date:       "123/123",
//	//	Cvv:        "456",
//	//	Meta:       " допом инфа",
//	//	UpdatedAt:  time.Now(),
//	//}
//	//
//	//err = b.db.Save(&user)
//	//if err != nil {
//	//	logger.Error()
//	//}
//	//
//	//user.ID++
//	//err = b.db.Save(&user)
//	//if err != nil {
//	//	logger.Error(err)
//	//}
//	//
//	//var cards []model.DataCard
//	//err = b.db.All(&cards)
//	//if err != nil {
//	//	logger.Error(err)
//	//}
//	//
//	//for i, v := range cards {
//	//	logger.Info(i, v)
//	//}
//	//
//	//b.db.Set("session", "754-3010", "ewsdfsdfsdfwrdff")
//	//
//	//var u string
//	//b.db.Get("session", "754-3010", &u)
//	//
//	//logger.Info(u)
//
//	return nil
//}
