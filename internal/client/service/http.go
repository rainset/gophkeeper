package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rainset/gophkeeper/internal/client/config"
	"github.com/rainset/gophkeeper/internal/client/model"
	server_model "github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/pkg/logger"
	"io"
	"net/http"
)

type HttpService struct {
	cfg         *config.Config
	AccessToken string
}

func NewHttpService(cfg *config.Config) *HttpService {
	return &HttpService{
		cfg: cfg,
	}
}

func (s *HttpService) SignIn(user model.User) (tokens model.Tokens, err error) {

	client := resty.New()

	res, err := client.R().
		SetBody(user).
		SetResult(&tokens).
		Post(s.cfg.ServerProtocol + "://" + s.cfg.ServerAddress + "/sign-in")

	if err != nil {
		logger.Error(err)
		return tokens, err
	}

	switch res.StatusCode() {

	case http.StatusOK:
		return tokens, err
		break
	case http.StatusUnauthorized:
		return tokens, ErrStatusUnauthorized
		break
	default:
		return tokens, err
		break
	}

	return tokens, err
}

func (s *HttpService) SignUp(user model.User) (tokens model.Tokens, err error) {

	client := resty.New()
	res, err := client.R().
		SetBody(user).
		SetResult(&tokens).
		Post(s.cfg.ServerProtocol + "://" + s.cfg.ServerAddress + "/sign-up")

	if err != nil {
		logger.Error(err)
		return tokens, err
	}

	switch res.StatusCode() {

	case http.StatusOK:
		return tokens, err
		break
	case http.StatusConflict:
		return tokens, ErrStatusLoginExists
		break
	default:
		return tokens, err
		break
	}

	return tokens, err
}

func (s *HttpService) PostRefreshToken(refreshToken string) (tokens model.Tokens, err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/refresh-token")
	rt := server_model.Tokens{RefreshToken: refreshToken}

	client := resty.New()
	_, err = client.R().
		SetBody(rt).
		SetResult(&tokens).
		Post(url)

	return tokens, err
}

func (s *HttpService) GetCardList(accessToken string) (items []model.DataCard, err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/card/list")

	client := resty.New()
	client.SetAuthToken(accessToken)

	resp, err := client.R().
		SetResult(&items).
		Get(url)

	if err != nil {
		logger.Info("url: ", url)
		logger.Info("statusCode: ", resp.StatusCode())
		logger.Info("resp: ", resp)
	}

	return items, err
}

func (s *HttpService) GetCredList(accessToken string) (items []model.DataCred, err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/cred/list")

	client := resty.New()
	client.SetAuthToken(accessToken)

	resp, err := client.R().
		SetResult(&items).
		Get(url)

	if err != nil {
		logger.Info("url: ", url)
		logger.Info("statusCode: ", resp.StatusCode())
		logger.Info("resp: ", resp)
	}

	return items, err
}

func (s *HttpService) GetTextList(accessToken string) (items []model.DataText, err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/text/list")

	client := resty.New()
	client.SetAuthToken(accessToken)

	resp, err := client.R().
		SetResult(&items).
		Get(url)

	if err != nil {
		logger.Info("url: ", url)
		logger.Info("statusCode: ", resp.StatusCode())
		logger.Info("resp: ", resp)
	}

	return items, err
}

func (s *HttpService) DeleteCard(accessToken string, extID int) (err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/card")

	var card server_model.DataCard
	card.ID = extID

	client := resty.New()
	client.SetAuthToken(accessToken)
	_, err = client.R().SetBody(card).Delete(url)
	if err != nil {
		logger.Error("DeleteCard: ", err)
	}

	return err
}

func (s *HttpService) DeleteCred(accessToken string, extID int) (err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/cred")

	var cred server_model.DataCred
	cred.ID = extID

	client := resty.New()
	client.SetAuthToken(accessToken)
	resp, err := client.R().SetBody(cred).Delete(url)
	if err != nil {
		logger.Error("DeleteCred: ", err)
	}

	logger.Error("DeleteCred")
	logger.Error(extID)
	logger.Error(accessToken)
	logger.Error(resp.StatusCode())
	logger.Error("")

	return err
}

func (s *HttpService) DeleteText(accessToken string, extID int) (err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/text")

	var card server_model.DataText
	card.ID = extID

	client := resty.New()
	client.SetAuthToken(accessToken)
	_, err = client.R().SetBody(card).Delete(url)
	if err != nil {
		logger.Error("DeleteText: ", err)
	}

	return err
}

func (s *HttpService) DeleteFile(accessToken string, extID int) (err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/file")

	var file server_model.DataFile
	file.ID = extID

	client := resty.New()
	client.SetAuthToken(accessToken)
	resp, err := client.R().SetBody(file).Delete(url)
	if err != nil {
		logger.Error("DeleteFile: ", err)
	}

	logger.Error("DeleteFile")
	logger.Error(extID)
	logger.Error(accessToken)
	logger.Error(resp.StatusCode())
	logger.Error("")

	return err
}

func (s *HttpService) GetFileList(accessToken string) (items []model.DataFile, err error) {

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/file/list")

	client := resty.New()
	client.SetAuthToken(accessToken)

	resp, err := client.R().
		SetResult(&items).
		Get(url)

	if err != nil {
		logger.Info("url: ", url)
		logger.Info("statusCode: ", resp.StatusCode())
		logger.Info("resp: ", resp)
	}

	return items, err
}

func (s *HttpService) DownloadFile(filePath string) (r io.ReadCloser, err error) {

	url := fmt.Sprintf("%s://%s/%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, filePath)
	resp, err := http.Get(url)
	if err != nil {
		logger.Error("http.Get: ", err)
	}
	return resp.Body, err
}
