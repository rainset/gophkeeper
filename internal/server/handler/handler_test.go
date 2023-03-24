package handler

import (
	"bytes"
	"context"
	"errors"
	"github.com/rainset/gophkeeper/internal/server/config"
	"github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/internal/server/service"
	"github.com/rainset/gophkeeper/internal/server/storage"
	"github.com/rainset/gophkeeper/internal/server/storage/file"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http/httptest"
	"testing"
)

func testUser() (tokens model.Tokens, err error) {

	cfg, err := config.ReadConfig()
	if err != nil {
		return tokens, err
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		return tokens, err
	}
	newService := service.New(store, storeFile, cfg)

	user := model.User{
		Login:    "test_handler_user_000000000",
		Password: "test_handler_user_000000000",
	}

	tokens, err = newService.SignIn(ctx, user)
	if err != nil {
		tokensSignUp, errSignUp := newService.SignUp(ctx, user)
		return tokensSignUp, errSignUp
	}

	return tokens, err
}

func TestHandler_DeleteCard(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		log.Fatal(err)
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "https://"+cfg.ServerAddress+"/store/card", bytes.NewBuffer([]byte(`{"id":0}`)))
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 401, w.Code)

}

func TestHandler_DeleteCred(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		log.Fatal(err)
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "https://"+cfg.ServerAddress+"/store/cred", bytes.NewBuffer([]byte(`{"id":0}`)))
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 401, w.Code)
}

func TestHandler_DeleteFile(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		log.Fatal(err)
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "https://"+cfg.ServerAddress+"/store/file", bytes.NewBuffer([]byte(`{"id":0}`)))
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 401, w.Code)
}

func TestHandler_DeleteText(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		log.Fatal(err)
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "https://"+cfg.ServerAddress+"/store/text", bytes.NewBuffer([]byte(`{"id":0}`)))
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 401, w.Code)
}

func TestHandler_FindAllCards(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://"+cfg.ServerAddress+"/store/card/list", nil)
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if w.Code != 200 && w.Code != 204 {
		// проверяем код ответа
		assert.Error(t, errors.New("response code only 200/204"))
	}
}

func TestHandler_FindAllCreds(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://"+cfg.ServerAddress+"/store/cred/list", nil)
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if w.Code != 200 && w.Code != 204 {
		// проверяем код ответа
		assert.Error(t, errors.New("response code only 200/204"))
	}
}

func TestHandler_FindAllFiles(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://"+cfg.ServerAddress+"/store/file/list", nil)
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if w.Code != 200 && w.Code != 204 {
		// проверяем код ответа
		assert.Error(t, errors.New("response code only 200/204"))
	}
}

func TestHandler_FindAllTexts(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://"+cfg.ServerAddress+"/store/text/list", nil)
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if w.Code != 200 && w.Code != 204 {
		// проверяем код ответа
		assert.Error(t, errors.New("response code only 200/204"))
	}
}

func TestHandler_FindCard(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/store/card", bytes.NewBuffer([]byte(`{"id":0}`)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 400, w.Code)
}

func TestHandler_FindCred(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/store/cred", bytes.NewBuffer([]byte(`{"id":0}`)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 400, w.Code)
}

func TestHandler_FindFile(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/store/file", bytes.NewBuffer([]byte(`{"id":0}`)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 400, w.Code)
}

func TestHandler_FindText(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/store/text", bytes.NewBuffer([]byte(`{"id":0}`)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 400, w.Code)
}

func TestHandler_Ping(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "https://"+cfg.ServerAddress+"/ping", nil)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 200, w.Code)
}

func TestHandler_RefreshToken(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/refresh-token", bytes.NewBuffer([]byte(`{"refresh_token":"`+tokens.RefreshToken+`"}`)))
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 200, w.Code)
}

func TestHandler_SaveCard(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()

	body := `{
	"title":"card",
    "number":"8977 4765 3453 9099",
    "date":"01/32",
    "cvv":"111",
    "meta":"meta"}`

	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/store/card", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 201, w.Code)
}

func TestHandler_SaveCred(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()

	body := `{
	"title":"credentials",
    "username":"username",
    "password":"password",
    "meta":"meta"}`

	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/store/cred", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 201, w.Code)
}

func TestHandler_SaveText(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()

	body := `{
	"title":"text",
    "text":"text",
    "meta":"meta"}`

	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/store/text", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 201, w.Code)
}

func TestHandler_SignIn(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()

	body := `{"login":"test_handler_user_000000000","password":"test_handler_user_000000000"}`

	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/sign-in", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 200, w.Code)
}

func TestHandler_SignKey(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()

	body := `{"login":"test_handler_user_000000000","password":"test_handler_user_000000000"}`

	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/sign-key", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 200, w.Code)
}

func TestHandler_SignUp(t *testing.T) {
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
		return
	}
	ctx := context.Background()
	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFile, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
		return
	}

	newService := service.New(store, storeFile, cfg)
	newHandler := NewHandler(newService)

	tokens, err := testUser()
	if err != nil {
		t.Error(err)
		return
	}

	if err != nil {
		t.Error(err)
		return
	}

	r := newHandler.Init()
	w := httptest.NewRecorder()

	body := `{"login":"test_handler_user_000000000","password":"test_handler_user_000000000"}`

	req := httptest.NewRequest("POST", "https://"+cfg.ServerAddress+"/sign-up", bytes.NewBuffer([]byte(body)))
	req.Header.Add("Authorization", "Bearer "+tokens.AccessToken)
	r.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	// проверяем код ответа
	assert.Equal(t, 400, w.Code)
}
