package service

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/rainset/gophkeeper/internal/client/config"
	"github.com/rainset/gophkeeper/internal/client/model"
	smodel "github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/pkg/logger"
)

type HTTPService struct {
	cfg *config.Config
}

func NewHTTPService(cfg *config.Config) *HTTPService {
	return &HTTPService{
		cfg: cfg,
	}
}

func (s *HTTPService) SignIn(user model.User) (tokens model.Tokens, err error) {
	client := resty.New()

	res, err := client.R().
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
	client := resty.New()
	res, err := client.R().
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

	client := resty.New()
	res, err := client.R().
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

func (s *HTTPService) GetCardList(accessToken string) (items []model.DataCard, err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/card/list")

	client := resty.New()
	client.SetAuthToken(accessToken)

	res, err := client.R().
		SetResult(&items).
		Get(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return items, ErrStatusUnauthorized
	default:
		return items, err
	}
}

func (s *HTTPService) GetCredList(accessToken string) (items []model.DataCred, err error) {
	url := fmt.Sprintf("%s://%s%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, "/store/cred/list")

	client := resty.New()
	client.SetAuthToken(accessToken)

	res, err := client.R().
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

	client := resty.New()
	client.SetAuthToken(accessToken)

	res, err := client.R().
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

	client := resty.New()
	client.SetAuthToken(accessToken)

	res, err := client.R().
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

	client := resty.New()
	client.SetAuthToken(accessToken)
	res, err := client.R().SetBody(card).Delete(url)

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

	client := resty.New()
	client.SetAuthToken(accessToken)
	res, err := client.R().SetBody(cred).Delete(url)

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

	client := resty.New()
	client.SetAuthToken(accessToken)
	res, err := client.R().SetBody(text).Delete(url)

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

	client := resty.New()
	client.SetAuthToken(accessToken)
	res, err := client.R().SetBody(file).Delete(url)

	switch res.StatusCode() {
	case http.StatusUnauthorized:
		return ErrStatusUnauthorized
	default:
		return err
	}
}

func (s *HTTPService) DownloadFile(filePath string) (r io.ReadCloser, err error) {
	url := fmt.Sprintf("%s://%s/%s", s.cfg.ServerProtocol, s.cfg.ServerAddress, filePath)

	resp, err := http.Get(url)

	if err != nil {
		logger.Error("http.Get: ", err)
	}

	return resp.Body, err
}
