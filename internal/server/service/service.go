package service

import (
	"errors"
	"github.com/rainset/gophkeeper/internal/server/config"
	"github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/internal/server/storage"
	"github.com/rainset/gophkeeper/internal/server/storage/file"
	"github.com/rainset/gophkeeper/pkg/auth"
	"github.com/rainset/gophkeeper/pkg/logger"
	"strconv"
	"time"
)

type Service struct {
	Store        storage.Interface
	StoreFiles   file.StorageFiles
	Cfg          *config.Config
	TokenManager auth.TokenManager
}

func New(store storage.Interface, storeFiles file.StorageFiles, cfg *config.Config) *Service {

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

func (s *Service) CreateSession(userId int) (model.Tokens, error) {
	var (
		res model.Tokens
		err error
	)

	res.AccessToken, err = s.TokenManager.NewJWT(strconv.Itoa(userId), s.Cfg.JWTAccessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.TokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	err = s.Store.SetRefreshToken(model.RefreshToken{UserID: userId, Token: res.RefreshToken, ExpiredAt: time.Now().Add(s.Cfg.JWTRefreshTokenTTL)})

	return res, err
}
func (s *Service) SignUp(user model.User) (tokens model.Tokens, err error) {

	userID, err := s.Store.CreateUser(user)
	if err != nil {
		return tokens, err
	}

	return s.CreateSession(userID)
}
func (s *Service) SignIn(user model.User) (tokens model.Tokens, err error) {

	userID, err := s.Store.GetUserIDByCredentials(user.Login, user.Password)

	if err != nil {
		return tokens, err
	}

	return s.CreateSession(userID)
}

func (s *Service) GetRefreshToken(token string) (tokens model.Tokens, err error) {

	userID, err := s.Store.GetRefreshTokenUserID(token)
	if err != nil {
		return tokens, err
	}
	if userID == 0 {
		return tokens, errors.New("refresh token is invalid")
	}

	return s.CreateSession(userID)
}

func (s *Service) ClearExpiredRefreshTokens() error {

	err := s.Store.ClearExpiredRefreshTokens()

	if err != nil {
		logger.Error("ClearExpiredRefreshTokens", err)
		return err
	}
	return err
}

func (s *Service) SaveCard(card model.DataCard) (err error) {
	err = card.Validate()
	if err != nil {
		logger.Error("SaveCard Validate()", err)
		return err
	}
	return s.Store.SaveCard(card)
}
func (s *Service) DeleteCard(ID, userID int) (err error) {
	return s.Store.DeleteCard(ID, userID)
}
func (s *Service) FindCard(ID, userID int) (card model.DataCard, err error) {
	return s.Store.FindCard(ID, userID)
}
func (s *Service) FindAllCards(userID int) (cards []model.DataCard, err error) {
	return s.Store.FindAllCards(userID)
}

func (s *Service) SaveFile(file model.DataFile) (err error) {
	err = file.Validate()
	if err != nil {
		logger.Error("SaveFile Validate()", err)
		return err
	}
	return s.Store.SaveFile(file)
}
func (s *Service) DeleteFile(fileID, userID int) (err error) {

	file, err := s.Store.FindFile(fileID, userID)
	if err != nil {
		logger.Error("find file to delete error", err)
		return err
	}

	if file.ID == 0 {
		logger.Error("file to delete not found", err)
		return errors.New("file to delete not found")
	}

	err = s.Store.DeleteFile(fileID, userID)
	if err != nil {
		logger.Error("find file to delete error", err)
		return err
	}

	err = s.StoreFiles.DeleteFile(file.Path)
	return err

}
func (s *Service) FindFile(ID, userID int) (file model.DataFile, err error) {
	return s.Store.FindFile(ID, userID)
}
func (s *Service) FindAllFiles(userID int) (files []model.DataFile, err error) {
	return s.Store.FindAllFiles(userID)
}

func (s *Service) SaveCred(cred model.DataCred) (err error) {
	err = cred.Validate()
	if err != nil {
		logger.Error("SaveCred Validate()", err)
		return err
	}
	return s.Store.SaveCred(cred)
}
func (s *Service) DeleteCred(ID, userID int) (err error) {
	return s.Store.DeleteCred(ID, userID)
}
func (s *Service) FindCred(ID, userID int) (cred model.DataCred, err error) {
	return s.Store.FindCred(ID, userID)
}
func (s *Service) FindAllCreds(userID int) (creds []model.DataCred, err error) {
	return s.Store.FindAllCreds(userID)
}

func (s *Service) SaveText(text model.DataText) (err error) {
	err = text.Validate()
	if err != nil {
		logger.Error("SaveText Validate()", err)
		return err
	}
	return s.Store.SaveText(text)
}
func (s *Service) DeleteText(ID, userID int) (err error) {
	return s.Store.DeleteText(ID, userID)
}
func (s *Service) FindText(ID, userID int) (text model.DataText, err error) {
	return s.Store.FindText(ID, userID)
}
func (s *Service) FindAllTexts(userID int) (texts []model.DataText, err error) {
	return s.Store.FindAllTexts(userID)
}
