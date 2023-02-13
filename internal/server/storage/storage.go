package storage

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rainset/gophkeeper/internal/server/model"
	"log"
	"time"
)

type Interface interface {
	CreateUser(user model.User) (userID int, err error)
	GetUserIDByCredentials(login, password string) (userID int, err error)

	SetRefreshToken(refreshToken model.RefreshToken) error
	GetRefreshTokenUserID(token string) (userID int, err error)
	ClearExpiredRefreshTokens() error

	SaveCard(card model.DataCard) (err error)
	FindCard(cardID, userID int) (card model.DataCard, err error)
	FindAllCards(userID int) (cards []model.DataCard, err error)
	DeleteCard(cardID, userID int) error

	SaveFile(file model.DataFile) (err error)
	DeleteFile(fileID, userID int) error
	FindFile(fileID, userID int) (file model.DataFile, err error)
	FindAllFiles(userID int) (files []model.DataFile, err error)

	SaveCred(cred model.DataCred) (err error)
	DeleteCred(credID, userID int) error
	FindCred(credID, userID int) (cred model.DataCred, err error)
	FindAllCreds(userID int) (creds []model.DataCred, err error)

	SaveText(text model.DataText) (err error)
	DeleteText(textID, userID int) error
	FindText(textID, userID int) (text model.DataText, err error)
	FindAllTexts(userID int) (texts []model.DataText, err error)
}

type Database struct {
	pgx *pgxpool.Pool
	ctx context.Context
}

func New(dataSourceName string) *Database {
	ctx := context.Background()
	db, err := pgxpool.New(context.Background(), dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("DB: connection initialized...")

	return &Database{
		pgx: db,
		ctx: ctx,
	}
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (d *Database) CreateUser(user model.User) (userID int, err error) {
	var hash = GetMD5Hash(user.Password)
	t := time.Now()
	sql := "INSERT INTO users (login,password,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING ID"
	err = d.pgx.QueryRow(d.ctx, sql, user.Login, hash, t, t).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return userID, ErrorUserAlreadyExists
			}
		}
	}
	return userID, err
}
func (d *Database) GetUserIDByCredentials(login, password string) (userID int, err error) {
	var hash = GetMD5Hash(password)
	var qPass string
	sql := "SELECT id,password FROM users WHERE login = $1"
	err = d.pgx.QueryRow(d.ctx, sql, login).Scan(&userID, &qPass)
	if err != nil {
		return userID, err
	}
	if hash != qPass {
		return userID, ErrorUserCredentials
	}

	return userID, err
}

func (d *Database) SetRefreshToken(in model.RefreshToken) error {

	sql := "INSERT INTO refresh_tokens (user_id, token , created_at, expired_at)  VALUES ($1, $2, $3, $4)"
	_, err := d.pgx.Exec(d.ctx, sql, in.UserID, in.Token, time.Now(), in.ExpiredAt)
	return err
}
func (d *Database) GetRefreshTokenUserID(token string) (userID int, err error) {

	sql := "SELECT user_id FROM refresh_tokens WHERE token=$1 AND expired_at>NOW()"
	err = d.pgx.QueryRow(d.ctx, sql, token).Scan(&userID)
	if err != nil {
		return userID, err
	}

	return userID, nil
}

func (d *Database) ClearExpiredRefreshTokens() error {

	sql := "DELETE FROM refresh_tokens WHERE expired_at < NOW()"
	_, err := d.pgx.Exec(d.ctx, sql)
	return err
}

func (d *Database) SaveCard(card model.DataCard) (err error) {
	if card.ID == 0 {
		sql := "INSERT INTO data_cards (user_id,title,number,date,cvv,meta,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7)"
		_, err = d.pgx.Exec(d.ctx, sql, card.UserID, card.Title, card.Number, card.Date, card.Cvv, card.Meta, time.Now())
	} else {
		sql := "UPDATE data_cards SET title=$1,number=$2,date=$3,cvv=$4,meta=$5,updated_at=$6 WHERE user_id=$7 AND id=$8"
		_, err = d.pgx.Exec(d.ctx, sql, card.Title, card.Number, card.Date, card.Cvv, card.Meta, time.Now(), card.UserID, card.ID)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return ErrorRowAlreadyExists
		}
	}

	return err
}
func (d *Database) FindCard(cardID, userID int) (card model.DataCard, err error) {
	sql := "SELECT id,title,number,date,cvv,meta,updated_at FROM data_cards WHERE id=$1 AND user_id = $2"
	err = pgxscan.Get(d.ctx, d.pgx, &card, sql, cardID, userID)
	return card, err
}
func (d *Database) FindAllCards(userID int) (cards []model.DataCard, err error) {
	sql := "SELECT id,title,number,date,cvv,meta,updated_at FROM data_cards WHERE user_id = $1 ORDER BY id DESC"
	err = pgxscan.Select(d.ctx, d.pgx, &cards, sql, userID)
	return cards, err
}
func (d *Database) DeleteCard(cardID, userID int) (err error) {
	sql := "DELETE FROM data_cards WHERE id=$1 AND user_id =$2"
	_, err = d.pgx.Exec(d.ctx, sql, cardID, userID)
	return err
}

func (d *Database) SaveFile(file model.DataFile) (err error) {
	if file.ID == 0 {
		sql := "INSERT INTO data_files (user_id,title,filename,path,meta,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
		_, err = d.pgx.Exec(d.ctx, sql, file.UserID, file.Title, file.Filename, file.Path, file.Meta, time.Now())
	} else {
		sql := "UPDATE data_files SET title=$1,filename=$2,path=$3,meta=$4, updated_at=$5 WHERE user_id=$6 AND id=$7"
		_, err = d.pgx.Exec(d.ctx, sql, file.Title, file.Filename, file.Path, file.Meta, time.Now(), file.UserID, file.ID)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return ErrorRowAlreadyExists
		}
	}

	return err
}
func (d *Database) DeleteFile(fileID, userID int) (err error) {
	sql := "DELETE  FROM data_files WHERE id=$1 AND user_id =$2"
	_, err = d.pgx.Exec(d.ctx, sql, fileID, userID)
	return err
}
func (d *Database) FindFile(fileID, userID int) (file model.DataFile, err error) {
	sql := "SELECT id,title,filename,path,meta,updated_at FROM data_files WHERE id=$1 AND user_id = $2"
	err = pgxscan.Get(d.ctx, d.pgx, &file, sql, fileID, userID)
	return file, err
}
func (d *Database) FindAllFiles(userID int) (files []model.DataFile, err error) {

	sql := "SELECT id,title,filename,path,meta,updated_at FROM data_files WHERE user_id = $1 ORDER BY id DESC"
	err = pgxscan.Select(d.ctx, d.pgx, &files, sql, userID)
	return files, err
}

func (d *Database) SaveCred(cred model.DataCred) (err error) {
	if cred.ID == 0 {
		sql := "INSERT INTO data_creds (user_id,title,username,password,meta,updated_at) VALUES ($1,$2,$3,$4,$5,$6)"
		_, err = d.pgx.Exec(d.ctx, sql, cred.UserID, cred.Title, cred.Username, cred.Password, cred.Meta, time.Now())
	} else {
		sql := "UPDATE data_creds SET title=$1,username=$2,password=$3,meta=$4, updated_at=$5 WHERE id=$6 AND user_id=$7"
		_, err = d.pgx.Exec(d.ctx, sql, cred.Title, cred.Username, cred.Password, cred.Meta, time.Now(), cred.ID, cred.UserID)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return ErrorRowAlreadyExists
		}
	}

	return err
}
func (d *Database) DeleteCred(credID, userID int) (err error) {
	sql := "DELETE FROM data_creds WHERE id=$1 AND user_id =$2"
	_, err = d.pgx.Exec(d.ctx, sql, credID, userID)
	return err
}
func (d *Database) FindCred(credID, userID int) (cred model.DataCred, err error) {
	sql := "SELECT id,title,username,password,meta,updated_at FROM data_creds WHERE id=$1 AND user_id = $2"
	err = pgxscan.Get(d.ctx, d.pgx, &cred, sql, credID, userID)
	return cred, err
}
func (d *Database) FindAllCreds(userID int) (creds []model.DataCred, err error) {

	sql := "SELECT id,title,username,password,meta,updated_at FROM data_creds WHERE user_id = $1 ORDER BY id DESC"
	err = pgxscan.Select(d.ctx, d.pgx, &creds, sql, userID)
	return creds, err
}

func (d *Database) SaveText(text model.DataText) (err error) {
	if text.ID == 0 {
		sql := "INSERT INTO data_text (user_id,title,text,meta,updated_at) VALUES ($1,$2,$3,$4,$5)"
		_, err = d.pgx.Exec(d.ctx, sql, text.UserID, text.Title, text.Text, text.Meta, time.Now())
	} else {
		sql := "UPDATE data_text SET title=$1,text=$2,meta=$3, updated_at=$4 WHERE id=$5 AND user_id=$6"
		_, err = d.pgx.Exec(d.ctx, sql, text.Title, text.Text, text.Meta, time.Now(), text.ID, text.UserID)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return ErrorRowAlreadyExists
		}
	}

	return err
}
func (d *Database) DeleteText(textID, userID int) (err error) {

	sql := "DELETE FROM data_text WHERE id=$1 AND user_id =$2"
	_, err = d.pgx.Exec(d.ctx, sql, textID, userID)
	return err
}
func (d *Database) FindText(textID, userID int) (text model.DataText, err error) {

	sql := "SELECT id,title,text,meta,updated_at FROM data_text WHERE id=$1 AND user_id = $2"
	err = pgxscan.Get(d.ctx, d.pgx, &text, sql, textID, userID)
	return text, err
}
func (d *Database) FindAllTexts(userID int) (texts []model.DataText, err error) {

	sql := "SELECT id,title,text,meta,updated_at FROM data_text WHERE user_id = $1 ORDER BY id DESC"
	err = pgxscan.Select(d.ctx, d.pgx, &texts, sql, userID)
	return texts, err
}
