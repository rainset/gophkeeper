package service

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rainset/gophkeeper/internal/client/config"
	"github.com/rainset/gophkeeper/internal/client/model"
	smodel "github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/pkg/logger"
)

type ResponseID struct {
	ID int `json:"id"`
}

type HTTPService struct {
	cfg    *config.Config
	client *resty.Client
}

func NewHTTPService(cfg *config.Config) *HTTPService {
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	return &HTTPService{
		cfg:    cfg,
		client: client,
	}
}

func (s *HTTPService) SignIn(user model.User) (tokens model.Tokens, err error) {
	res, err := s.client.R().
		SetBody(user).
		SetResult(&tokens).
		Post(s.cfg.ServerProtocol + "://" + s.cfg.ServerAddress + "/sign-in")

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return tokens, ErrStatusUnauthorized
	default:
		return tokens, err
	}
}

func (s *HTTPService) SignUp(user model.User) (tokens model.Tokens, err error) {
	res, err := s.client.R().
		SetBody(user).
		SetResult(&tokens).
		Post(s.cfg.ServerProtocol + "://" + s.cfg.ServerAddress + "/sign-up")

	switch res.StatusCode() {
	case http.StatusConflict:
		return tokens, ErrStatusLoginExists
	default:
		return tokens, err
	}
}

func (s *HTTPService) PostRefreshToken(refreshToken string) (tokens model.Tokens, err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/refresh-token")
	rt := smodel.Tokens{RefreshToken: refreshToken}
	res, err := s.client.R().
		SetBody(rt).
		SetResult(&tokens).
		Post(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return tokens, ErrStatusUnauthorized
	default:
		return tokens, err
	}
}

func (s *HTTPService) GetSignKey(accessToken string, login, password string) (signKey string, err error) {

	user := model.User{
		Login:    login,
		Password: password,
	}

	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/sign-key")

	var resp struct {
		SignKey string `json:"sign_key"`
	}
	s.client.SetAuthToken(accessToken)
	res, err := s.client.R().
		SetBody(user).
		SetResult(&resp).
		Post(url)

	signKey = resp.SignKey

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return signKey, ErrStatusUnauthorized
	default:
		return signKey, err
	}
}

func (s *HTTPService) GetCardList(accessToken string) (items []*model.DataCard, err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/card/list")

	s.client.SetAuthToken(accessToken)
	res, err := s.client.R().
		SetResult(&items).
		Get(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return items, ErrStatusUnauthorized
	default:
		return items, err
	}
}

func (s *HTTPService) GetCredList(accessToken string) (items []*model.DataCred, err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/cred/list")

	s.client.SetAuthToken(accessToken)

	res, err := s.client.R().
		SetResult(&items).
		Get(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return items, ErrStatusUnauthorized
	default:
		return items, err
	}
}

func (s *HTTPService) GetTextList(accessToken string) (items []model.DataText, err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/text/list")

	s.client.SetAuthToken(accessToken)

	res, err := s.client.R().
		SetResult(&items).
		Get(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return items, ErrStatusUnauthorized
	default:
		return items, err
	}
}

func (s *HTTPService) GetFileList(accessToken string) (items []model.DataFile, err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/file/list")

	s.client.SetAuthToken(accessToken)

	res, err := s.client.R().
		SetResult(&items).
		Get(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return items, ErrStatusUnauthorized
	default:
		return items, err
	}
}

func (s *HTTPService) DeleteCard(accessToken string, extID int) (err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/card")

	card := smodel.DataCard{ID: extID}
	s.client.SetAuthToken(accessToken)

	res, err := s.client.R().SetBody(card).Delete(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return ErrStatusUnauthorized
	default:
		return err
	}
}

func (s *HTTPService) DeleteCred(accessToken string, extID int) (err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/cred")

	cred := smodel.DataCred{ID: extID}

	s.client.SetAuthToken(accessToken)
	res, err := s.client.R().SetBody(cred).Delete(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return ErrStatusUnauthorized
	default:
		return err
	}
}

func (s *HTTPService) DeleteText(accessToken string, extID int) (err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/text")

	text := smodel.DataCred{ID: extID}

	s.client.SetAuthToken(accessToken)
	res, err := s.client.R().SetBody(text).Delete(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return ErrStatusUnauthorized
	default:
		return err
	}
}

func (s *HTTPService) DeleteFile(accessToken string, extID int) (err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/file")

	file := smodel.DataFile{ID: extID}

	s.client.SetAuthToken(accessToken)
	res, err := s.client.R().SetBody(file).Delete(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return ErrStatusUnauthorized
	default:
		return err
	}
}

func (s *HTTPService) DownloadFile(filePath string) (r io.ReadCloser, err error) {
	url := fmt.Sprintf("%s://%s/%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, filePath)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get(url)

	if err != nil {
		logger.Error("http.Get: ", err)
	}

	return resp.Body, err
}

func (s *HTTPService) AddCard(accessToken string, card smodel.DataCard) (id int, err error) {
	var rb ResponseID
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/card")
	s.client.SetAuthToken(accessToken)
	res, err := s.client.R().SetBody(card).SetResult(&rb).Post(url)

	logger.Info(string(res.Body()))
	logger.Info("rb ", rb.ID)
	logger.Info("cardID ", card.ID)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return rb.ID, ErrStatusUnauthorized
	default:
		return rb.ID, err
	}
}

func (s *HTTPService) AddCred(accessToken string, cred smodel.DataCred) (id int, err error) {
	var rb ResponseID
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/cred")

	s.client.SetAuthToken(accessToken)
	res, err := s.client.R().SetBody(cred).SetResult(rb).Post(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return rb.ID, ErrStatusUnauthorized
	default:
		return rb.ID, err
	}
}

func (s *HTTPService) AddText(accessToken string, text smodel.DataText) (id int, err error) {
	var rb ResponseID
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/text")

	s.client.SetAuthToken(accessToken)
	res, err := s.client.R().SetBody(text).SetResult(rb).Post(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return rb.ID, ErrStatusUnauthorized
	default:
		return rb.ID, err
	}
}

func (s *HTTPService) AddFile(accessToken string, file smodel.DataFile) (id int, err error) {
	var rb ResponseID
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/file")

	s.client.SetAuthToken(accessToken)

	// Multipart of form fields and files
	res, err := s.client.R().
		SetFiles(map[string]string{
			"file": file.Path,
		}).
		SetFormData(map[string]string{
			"title": file.Title,
			"meta":  file.Meta,
		}).SetResult(rb).Post(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return rb.ID, ErrStatusUnauthorized
	default:
		return rb.ID, err
	}
}
