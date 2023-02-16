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

var ErrUserNotInitialized = errors.New("user not initialized")

func (b *Base) SetUserConfig(config model.UserConfig) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	err = b.db.From(b.user).Set("store", "config", config)

	return err
}

func (b *Base) GetUserConfig() (c model.UserConfig, err error) {
	if b.user == "" {
		return c, ErrUserNotInitialized
	}

	err = b.db.From(b.user).Get("store", "config", &c)

	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return c, nil
		}
	}

	return c, err
}

func (b *Base) AddCard(card model.DataCard) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	err = b.db.From(b.user).Save(&card)
	if err != nil {
		logger.Error()
	}

	return err
}

func (b *Base) GetCard(localID int) (card model.DataCard, err error) {
	if b.user == "" {
		return card, ErrUserNotInitialized
	}

	err = b.db.From(b.user).One("LocalID", localID, &card)
	if err != nil {
		logger.Error()
	}

	return card, err
}

func (b *Base) GetAllCards() (cards []model.DataCard, err error) {
	if b.user == "" {
		return cards, ErrUserNotInitialized
	}

	err = b.db.From(b.user).All(&cards, storm.Reverse())
	if err != nil {
		logger.Error(err)
	}

	return cards, err
}

func (b *Base) DeleteCard(localID int) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	var card model.DataCard

	card.LocalID = localID
	err = b.db.From(b.user).DeleteStruct(&card)

	if err != nil {
		logger.Error(err)
	}

	return err
}

func (b *Base) AddCred(cred model.DataCred) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	err = b.db.From(b.user).Save(&cred)
	if err != nil {
		logger.Error()
	}

	return err
}

func (b *Base) GetCred(localID int) (cred model.DataCred, err error) {
	if b.user == "" {
		return cred, ErrUserNotInitialized
	}

	err = b.db.From(b.user).One("LocalID", localID, &cred)

	return cred, err
}

func (b *Base) GetAllCreds() (creds []model.DataCred, err error) {
	if b.user == "" {
		return creds, ErrUserNotInitialized
	}

	err = b.db.From(b.user).All(&creds, storm.Reverse())

	return creds, err
}

func (b *Base) DeleteCred(localID int) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	var cred model.DataCred
	cred.LocalID = localID
	err = b.db.From(b.user).DeleteStruct(&cred)

	return err
}

func (b *Base) AddText(text model.DataText) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	err = b.db.From(b.user).Save(&text)

	return err
}

func (b *Base) GetText(localID int) (text model.DataText, err error) {
	if b.user == "" {
		return text, ErrUserNotInitialized
	}

	err = b.db.From(b.user).One("LocalID", localID, &text)

	return text, err
}

func (b *Base) GetAllTexts() (texts []model.DataText, err error) {
	if b.user == "" {
		return texts, ErrUserNotInitialized
	}

	err = b.db.From(b.user).All(&texts, storm.Reverse())

	return texts, err
}

func (b *Base) DeleteText(localID int) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	var text model.DataText
	text.LocalID = localID
	err = b.db.From(b.user).DeleteStruct(&text)

	return err
}

func (b *Base) AddFile(file model.DataFile) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	err = b.db.From(b.user).Save(&file)

	return err
}

func (b *Base) GetFile(localID int) (file model.DataFile, err error) {
	if b.user == "" {
		return file, ErrUserNotInitialized
	}

	err = b.db.From(b.user).One("LocalID", localID, &file)

	return file, err
}

func (b *Base) GetAllFiles() (files []model.DataFile, err error) {
	if b.user == "" {
		return files, ErrUserNotInitialized
	}

	err = b.db.From(b.user).All(&files, storm.Reverse())
	if err != nil {
		logger.Error(err)
	}

	return files, err
}

func (b *Base) DeleteFile(localID int) (err error) {
	if b.user == "" {
		return ErrUserNotInitialized
	}

	var f model.DataFile
	f.LocalID = localID

	err = b.db.From(b.user).DeleteStruct(&f)
	if err != nil {
		logger.Error(err)
	}

	return err
}
