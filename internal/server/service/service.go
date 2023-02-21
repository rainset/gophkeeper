package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/rainset/gophkeeper/internal/server/config"
	"github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/internal/server/storage"
	"github.com/rainset/gophkeeper/internal/server/storage/file"
	"github.com/rainset/gophkeeper/pkg/auth"
	"github.com/rainset/gophkeeper/pkg/logger"
)

type Service struct {
	Store        storage.Interface
	StoreFiles   *file.StorageFiles
	Cfg          *config.Config
	TokenManager auth.TokenManager
}

func New(store storage.Interface, storeFiles *file.StorageFiles, cfg *config.Config) *Service {
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		logger.Error(err)
	}

	return &Service{
		Cfg:          cfg,
		Store:        store,
		StoreFiles:   storeFiles,
		TokenManager: tokenManager,
	}
}

func (s *Service) GetSignKey(ctx context.Context, login, password string) (signKey string, err error) {
	return s.Store.GetSignKey(ctx, login, password)
}

func (s *Service) ClearExpiredRefreshTokens(ctx context.Context) error {
	err := s.Store.ClearExpiredRefreshTokens(ctx)
	if err != nil {
		logger.Error("ClearExpiredRefreshTokens", err)

		return err
	}

	return err
}

func (s *Service) CreateSession(ctx context.Context, userID int) (model.Tokens, error) {
	var (
		res model.Tokens
		err error
	)

	accessTTL, err := time.ParseDuration(s.Cfg.JWTAccessTokenTTL)
	if err != nil {
		return res, err
	}

	res.AccessToken, err = s.TokenManager.NewJWT(strconv.Itoa(userID), accessTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.TokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	refreshTokenTTL, err := time.ParseDuration(s.Cfg.JWTRefreshTokenTTL)
	if err != nil {
		return res, err
	}

	err = s.Store.SetRefreshToken(ctx, model.RefreshToken{UserID: userID, Token: res.RefreshToken, ExpiredAt: time.Now().Add(refreshTokenTTL)})

	return res, err
}

func (s *Service) SignUp(ctx context.Context, user model.User) (tokens model.Tokens, err error) {
	userID, err := s.Store.CreateUser(ctx, user)
	if err != nil {
		return tokens, err
	}

	return s.CreateSession(ctx, userID)
}

func (s *Service) SignIn(ctx context.Context, user model.User) (tokens model.Tokens, err error) {
	userID, err := s.Store.GetUserIDByCredentials(ctx, user.Login, user.Password)
	if err != nil {
		return tokens, err
	}

	return s.CreateSession(ctx, userID)
}

func (s *Service) GetRefreshToken(ctx context.Context, token string) (tokens model.Tokens, err error) {
	userID, err := s.Store.GetRefreshTokenUserID(ctx, token)
	if err != nil {
		return tokens, err
	}

	if userID == 0 {
		return tokens, errors.New("refresh token is invalid")
	}

	return s.CreateSession(ctx, userID)
}

func (s *Service) SaveCard(ctx context.Context, card model.DataCard) (id int, err error) {
	err = card.Validate()
	if err != nil {
		return id, err
	}
	return s.Store.SaveCard(ctx, card)
}

func (s *Service) DeleteCard(ctx context.Context, cardID, userID int) (err error) {
	return s.Store.DeleteCard(ctx, cardID, userID)
}

func (s *Service) FindCard(ctx context.Context, cardID, userID int) (card model.DataCard, err error) {
	return s.Store.FindCard(ctx, cardID, userID)
}

func (s *Service) FindAllCards(ctx context.Context, userID int) (cards []model.DataCard, err error) {
	return s.Store.FindAllCards(ctx, userID)
}

func (s *Service) SaveFile(ctx context.Context, file model.DataFile) (err error) {
	err = file.Validate()
	if err != nil {
		logger.Error("SaveFile Validate()", err)

		return err
	}

	return s.Store.SaveFile(ctx, file)
}

func (s *Service) DeleteFile(ctx context.Context, fileID, userID int) (err error) {
	file, err := s.Store.FindFile(ctx, fileID, userID)
	if err != nil {
		logger.Error("find file to delete error", err)

		return err
	}

	if file.ID == 0 {
		logger.Error("file to delete not found", err)

		return errors.New("file to delete not found")
	}

	err = s.Store.DeleteFile(ctx, fileID, userID)
	if err != nil {
		logger.Error("find file to delete error", err)

		return err
	}

	err = s.StoreFiles.DeleteFile(file.Path)

	return err
}

func (s *Service) FindFile(ctx context.Context, fileID, userID int) (file model.DataFile, err error) {
	return s.Store.FindFile(ctx, fileID, userID)
}

func (s *Service) FindAllFiles(ctx context.Context, userID int) (files []model.DataFile, err error) {
	return s.Store.FindAllFiles(ctx, userID)
}

func (s *Service) SaveCred(ctx context.Context, cred model.DataCred) (err error) {
	err = cred.Validate()
	if err != nil {
		logger.Error("SaveCred Validate()", err)

		return err
	}

	return s.Store.SaveCred(ctx, cred)
}

func (s *Service) DeleteCred(ctx context.Context, credID, userID int) (err error) {
	return s.Store.DeleteCred(ctx, credID, userID)
}

func (s *Service) FindCred(ctx context.Context, credID, userID int) (cred model.DataCred, err error) {
	return s.Store.FindCred(ctx, credID, userID)
}

func (s *Service) FindAllCreds(ctx context.Context, userID int) (creds []model.DataCred, err error) {
	return s.Store.FindAllCreds(ctx, userID)
}

func (s *Service) SaveText(ctx context.Context, text model.DataText) (err error) {
	err = text.Validate()
	if err != nil {
		logger.Error("SaveText Validate()", err)

		return err
	}

	return s.Store.SaveText(ctx, text)
}

func (s *Service) DeleteText(ctx context.Context, textID, userID int) (err error) {
	return s.Store.DeleteText(ctx, textID, userID)
}

func (s *Service) FindText(ctx context.Context, textID, userID int) (text model.DataText, err error) {
	return s.Store.FindText(ctx, textID, userID)
}

func (s *Service) FindAllTexts(ctx context.Context, userID int) (texts []model.DataText, err error) {
	return s.Store.FindAllTexts(ctx, userID)
}
