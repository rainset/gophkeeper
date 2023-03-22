package storage

import (
	"context"
	"crypto/aes"
	"errors"
	"fmt"
	"github.com/rainset/gophkeeper/pkg/crypt"
	"github.com/rainset/gophkeeper/pkg/hash"
	"log"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rainset/gophkeeper/internal/server/model"
)

type Interface interface {
	CreateUser(ctx context.Context, user model.User) (userID int, err error)
	GetUserIDByCredentials(ctx context.Context, login, password string) (userID int, err error)
	GetSignKey(ctx context.Context, login, password string) (signKey string, err error)

	SetRefreshToken(ctx context.Context, refreshToken model.RefreshToken) error
	GetRefreshTokenUserID(ctx context.Context, token string) (userID int, err error)
	ClearExpiredRefreshTokens(ctx context.Context) error

	SaveCard(ctx context.Context, card model.DataCard) (id int, err error)
	FindCard(ctx context.Context, cardID, userID int) (card model.DataCard, err error)
	FindAllCards(ctx context.Context, userID int) (cards []model.DataCard, err error)
	DeleteCard(ctx context.Context, cardID, userID int) error

	SaveFile(ctx context.Context, file model.DataFile) (id int, err error)
	DeleteFile(ctx context.Context, fileID, userID int) error
	FindFile(ctx context.Context, fileID, userID int) (file model.DataFile, err error)
	FindAllFiles(ctx context.Context, userID int) (files []model.DataFile, err error)

	SaveCred(ctx context.Context, cred model.DataCred) (id int, err error)
	DeleteCred(ctx context.Context, credID, userID int) error
	FindCred(ctx context.Context, credID, userID int) (cred model.DataCred, err error)
	FindAllCreds(ctx context.Context, userID int) (creds []model.DataCred, err error)

	SaveText(ctx context.Context, text model.DataText) (id int, err error)
	DeleteText(ctx context.Context, textID, userID int) error
	FindText(ctx context.Context, textID, userID int) (text model.DataText, err error)
	FindAllTexts(ctx context.Context, userID int) (texts []model.DataText, err error)
}

type Database struct {
	pgx *pgxpool.Pool
}

func New(ctx context.Context, dataSourceName string) *Database {

	db, err := pgxpool.New(ctx, dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	log.Print("DB: connection initialized...")

	return &Database{
		pgx: db,
	}
}

func (d *Database) Close() {
	d.pgx.Close()
}

func (d *Database) GetSignKey(ctx context.Context, login, password string) (signKey string, err error) {
	var passHash = hash.Md5(password)

	sql := "SELECT sign_key FROM users WHERE login = $1 AND password = $2"
	err = d.pgx.QueryRow(ctx, sql, login, passHash).Scan(&signKey)
	return signKey, fmt.Errorf("db.GetSignKey: %w", err)
}

func (d *Database) CreateUser(ctx context.Context, user model.User) (userID int, err error) {
	var hMD5 = hash.Md5(user.Password)
	t := time.Now()

	h, err := hash.GenerateRandom(2 * aes.BlockSize)
	if err != nil {
		return userID, fmt.Errorf("db.CreateUser: %w", err)
	}

	signKey := crypt.EncodeBase64(h)

	sql := "INSERT INTO users (login,password,sign_key,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id"

	err = d.pgx.QueryRow(ctx, sql, user.Login, hMD5, signKey, t, t).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return userID, ErrorUserAlreadyExists
			}
		}
	}

	return userID, fmt.Errorf("db.CreateUser: %w", err)
}
func (d *Database) GetUserIDByCredentials(ctx context.Context, login, password string) (userID int, err error) {
	var qPass string
	var hash = hash.Md5(password)

	sql := "SELECT id,password FROM users WHERE login = $1"

	err = d.pgx.QueryRow(ctx, sql, login).Scan(&userID, &qPass)
	if err != nil {
		return userID, fmt.Errorf("db.GetUserIDByCredentials: %w", err)
	}

	if hash != qPass {
		return userID, ErrorUserCredentials
	}

	return userID, fmt.Errorf("db.GetUserIDByCredentials: %w", err)
}

func (d *Database) SetRefreshToken(ctx context.Context, in model.RefreshToken) error {
	sql := "INSERT INTO refresh_tokens (user_id, token , created_at, expired_at)  VALUES ($1, $2, $3, $4)"
	_, err := d.pgx.Exec(ctx, sql, in.UserID, in.Token, time.Now(), in.ExpiredAt)

	return fmt.Errorf("db.SetRefreshToken: %w", err)
}
func (d *Database) GetRefreshTokenUserID(ctx context.Context, token string) (userID int, err error) {
	sql := "SELECT user_id FROM refresh_tokens WHERE token=$1 AND expired_at>NOW()"

	err = d.pgx.QueryRow(ctx, sql, token).Scan(&userID)
	if err != nil {
		return userID, fmt.Errorf("db.GetRefreshTokenUserID: %w", err)
	}

	return userID, nil
}

func (d *Database) ClearExpiredRefreshTokens(ctx context.Context) error {
	sql := "DELETE FROM refresh_tokens WHERE expired_at < NOW()"
	_, err := d.pgx.Exec(ctx, sql)

	return fmt.Errorf("db.ClearExpiredRefreshTokens: %w", err)
}

func (d *Database) SaveCard(ctx context.Context, card model.DataCard) (id int, err error) {
	if card.ID == 0 {
		sql := "INSERT INTO data_cards (user_id,title,number,date,cvv,meta,updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id"
		err = d.pgx.QueryRow(ctx, sql, card.UserID, card.Title, card.Number, card.Date, card.Cvv, card.Meta, card.UpdatedAt).Scan(&id)
	} else {
		id = card.ID
		sql := "UPDATE data_cards SET title=$1,number=$2,date=$3,cvv=$4,meta=$5,updated_at=$6 WHERE user_id=$7 AND id=$8"
		_, err = d.pgx.Exec(ctx, sql, card.Title, card.Number, card.Date, card.Cvv, card.Meta, card.UpdatedAt, card.UserID, card.ID)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return id, ErrorRowAlreadyExists
		}
	}

	return id, fmt.Errorf("db.SaveCard: %w", err)
}
func (d *Database) FindCard(ctx context.Context, cardID, userID int) (card model.DataCard, err error) {
	sql := "SELECT id,title,number,date,cvv,meta,updated_at FROM data_cards WHERE id=$1 AND user_id = $2"
	err = pgxscan.Get(ctx, d.pgx, &card, sql, cardID, userID)

	return card, fmt.Errorf("db.FindCard: %w", err)
}
func (d *Database) FindAllCards(ctx context.Context, userID int) (cards []model.DataCard, err error) {
	sql := "SELECT id,title,number,date,cvv,meta,updated_at FROM data_cards WHERE user_id = $1 ORDER BY id DESC"
	err = pgxscan.Select(ctx, d.pgx, &cards, sql, userID)

	return cards, fmt.Errorf("db.FindAllCards: %w", err)
}
func (d *Database) DeleteCard(ctx context.Context, cardID, userID int) (err error) {
	sql := "DELETE FROM data_cards WHERE id=$1 AND user_id =$2"
	_, err = d.pgx.Exec(ctx, sql, cardID, userID)

	return fmt.Errorf("db.DeleteCard: %w", err)
}

func (d *Database) SaveFile(ctx context.Context, file model.DataFile) (id int, err error) {
	if file.ID == 0 {
		sql := "INSERT INTO data_files (user_id,title,filename,path,meta,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
		err = d.pgx.QueryRow(ctx, sql, file.UserID, file.Title, file.Filename, file.Path, file.Meta, file.UpdatedAt).Scan(&id)
	} else {
		id = file.ID
		sql := "UPDATE data_files SET title=$1,filename=$2,path=$3,meta=$4, updated_at=$5 WHERE user_id=$6 AND id=$7"
		_, err = d.pgx.Exec(ctx, sql, file.Title, file.Filename, file.Path, file.Meta, file.UpdatedAt, file.UserID, file.ID)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return id, ErrorRowAlreadyExists
		}
	}

	return id, fmt.Errorf("db.SaveFile: %w", err)
}
func (d *Database) DeleteFile(ctx context.Context, fileID, userID int) (err error) {
	sql := "DELETE  FROM data_files WHERE id=$1 AND user_id =$2"
	_, err = d.pgx.Exec(ctx, sql, fileID, userID)

	return fmt.Errorf("db.DeleteFile: %w", err)
}
func (d *Database) FindFile(ctx context.Context, fileID, userID int) (file model.DataFile, err error) {
	sql := "SELECT id,title,filename,path,meta,updated_at FROM data_files WHERE id=$1 AND user_id = $2"
	err = pgxscan.Get(ctx, d.pgx, &file, sql, fileID, userID)

	return file, fmt.Errorf("db.FindFile: %w", err)
}
func (d *Database) FindAllFiles(ctx context.Context, userID int) (files []model.DataFile, err error) {
	sql := "SELECT id,title,filename,path,meta,updated_at FROM data_files WHERE user_id = $1 ORDER BY id DESC"
	err = pgxscan.Select(ctx, d.pgx, &files, sql, userID)

	return files, fmt.Errorf("db.FindAllFiles: %w", err)
}

func (d *Database) SaveCred(ctx context.Context, cred model.DataCred) (id int, err error) {
	if cred.ID == 0 {
		sql := "INSERT INTO data_creds (user_id,title,username,password,meta,updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id"
		err = d.pgx.QueryRow(ctx, sql, cred.UserID, cred.Title, cred.Username, cred.Password, cred.Meta, cred.UpdatedAt).Scan(&id)
	} else {
		id = cred.ID
		sql := "UPDATE data_creds SET title=$1,username=$2,password=$3,meta=$4, updated_at=$5 WHERE id=$6 AND user_id=$7"
		_, err = d.pgx.Exec(ctx, sql, cred.Title, cred.Username, cred.Password, cred.Meta, cred.UpdatedAt, cred.ID, cred.UserID)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return id, ErrorRowAlreadyExists
		}
	}

	return id, fmt.Errorf("db.SaveCred: %w", err)
}
func (d *Database) DeleteCred(ctx context.Context, credID, userID int) (err error) {
	sql := "DELETE FROM data_creds WHERE id=$1 AND user_id =$2"
	_, err = d.pgx.Exec(ctx, sql, credID, userID)

	return fmt.Errorf("db.DeleteCred: %w", err)
}
func (d *Database) FindCred(ctx context.Context, credID, userID int) (cred model.DataCred, err error) {
	sql := "SELECT id,title,username,password,meta,updated_at FROM data_creds WHERE id=$1 AND user_id = $2"
	err = pgxscan.Get(ctx, d.pgx, &cred, sql, credID, userID)

	return cred, fmt.Errorf("db.FindCred: %w", err)
}
func (d *Database) FindAllCreds(ctx context.Context, userID int) (creds []model.DataCred, err error) {
	sql := "SELECT id,title,username,password,meta,updated_at FROM data_creds WHERE user_id = $1 ORDER BY id DESC"
	err = pgxscan.Select(ctx, d.pgx, &creds, sql, userID)

	return creds, fmt.Errorf("db.FindAllCreds: %w", err)
}

func (d *Database) SaveText(ctx context.Context, text model.DataText) (id int, err error) {
	if text.ID == 0 {
		sql := "INSERT INTO data_text (user_id,title,text,meta,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id"
		err = d.pgx.QueryRow(ctx, sql, text.UserID, text.Title, text.Text, text.Meta, text.UpdatedAt).Scan(&id)
	} else {
		id = text.ID
		sql := "UPDATE data_text SET title=$1,text=$2,meta=$3, updated_at=$4 WHERE id=$5 AND user_id=$6"
		_, err = d.pgx.Exec(ctx, sql, text.Title, text.Text, text.Meta, text.UpdatedAt, text.ID, text.UserID)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return id, ErrorRowAlreadyExists
		}
	}

	return id, fmt.Errorf("db.SaveText: %w", err)
}
func (d *Database) DeleteText(ctx context.Context, textID, userID int) (err error) {
	sql := "DELETE FROM data_text WHERE id=$1 AND user_id =$2"
	_, err = d.pgx.Exec(ctx, sql, textID, userID)

	return fmt.Errorf("db.DeleteText: %w", err)
}
func (d *Database) FindText(ctx context.Context, textID, userID int) (text model.DataText, err error) {
	sql := "SELECT id,title,text,meta,updated_at FROM data_text WHERE id=$1 AND user_id = $2"
	err = pgxscan.Get(ctx, d.pgx, &text, sql, textID, userID)

	return text, fmt.Errorf("db.FindText: %w", err)
}
func (d *Database) FindAllTexts(ctx context.Context, userID int) (texts []model.DataText, err error) {
	sql := "SELECT id,title,text,meta,updated_at FROM data_text WHERE user_id = $1 ORDER BY id DESC"
	err = pgxscan.Select(ctx, d.pgx, &texts, sql, userID)

	return texts, fmt.Errorf("db.FindAllTexts: %w", err)
}
