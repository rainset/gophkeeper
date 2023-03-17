package service

import (
	"context"
	"github.com/rainset/gophkeeper/internal/server/config"
	"github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/internal/server/storage"
	"github.com/rainset/gophkeeper/internal/server/storage/file"
	"github.com/rainset/gophkeeper/pkg/auth"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {

	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}

	ctx := context.Background()

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type args struct {
		store      storage.Interface
		storeFiles *file.StorageFiles
		cfg        *config.Config
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "new service",
			args: args{
				store:      store,
				storeFiles: storeFiles,
				cfg:        cfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.store, tt.args.storeFiles, tt.args.cfg)
			assert.NotNil(t, got)
		})
	}
}

func TestService_ClearExpiredRefreshTokens(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ClearExpiredRefreshTokens",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{ctx: ctx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			err = s.ClearExpiredRefreshTokens(tt.args.ctx)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestService_CreateSession(t *testing.T) {

	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.Tokens
		wantErr bool
	}{
		{
			name: "create session",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{ctx: ctx, userID: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			_, err = s.CreateSession(tt.args.ctx, tt.args.userID)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestService_DeleteCard(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		cardID int
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete card",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{ctx: ctx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			err = s.DeleteCard(tt.args.ctx, tt.args.cardID, tt.args.userID)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestService_DeleteCred(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		credID int
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete cred",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{ctx: ctx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			err = s.DeleteCred(tt.args.ctx, tt.args.credID, tt.args.userID)
			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestService_DeleteFile(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		fileID int
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete file",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				fileID: 0,
				userID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			err = s.DeleteFile(tt.args.ctx, tt.args.fileID, tt.args.userID)
			if err != nil {
				if err.Error() == "scanning one: no rows in result set" {
					return
				}
				t.Error(err)
			}
		})
	}
}

func TestService_DeleteText(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		textID int
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "delete text",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				textID: 0,
				userID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			if err := s.DeleteText(tt.args.ctx, tt.args.textID, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_FindAllCards(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCards []model.DataCard
		wantErr   bool
	}{
		{
			name: "find all cards",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				userID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotCards, err := s.FindAllCards(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllCards() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCards, tt.wantCards) {
				t.Errorf("FindAllCards() gotCards = %v, want %v", gotCards, tt.wantCards)
			}
		})
	}
}

func TestService_FindAllCreds(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCreds []model.DataCred
		wantErr   bool
	}{
		{
			name: "find all creds",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				userID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotCreds, err := s.FindAllCreds(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllCreds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCreds, tt.wantCreds) {
				t.Errorf("FindAllCreds() gotCreds = %v, want %v", gotCreds, tt.wantCreds)
			}
		})
	}
}

func TestService_FindAllFiles(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantFiles []model.DataFile
		wantErr   bool
	}{
		{
			name: "find all files",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				userID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotFiles, err := s.FindAllFiles(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFiles, tt.wantFiles) {
				t.Errorf("FindAllFiles() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
			}
		})
	}
}

func TestService_FindAllTexts(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		userID int
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantTexts []model.DataText
		wantErr   bool
	}{
		{
			name: "find all texts",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				userID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotTexts, err := s.FindAllTexts(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAllTexts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTexts, tt.wantTexts) {
				t.Errorf("FindAllTexts() gotTexts = %v, want %v", gotTexts, tt.wantTexts)
			}
		})
	}
}

func TestService_FindCard(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		cardID int
		userID int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "find card",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				cardID: 0,
				userID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			_, err := s.FindCard(tt.args.ctx, tt.args.cardID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestService_FindCred(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		credID int
		userID int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCred model.DataCred
		wantErr  bool
	}{
		{
			name: "find cred",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				credID: 0,
				userID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			_, err := s.FindCred(tt.args.ctx, tt.args.credID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindCred() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_FindFile(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		fileID int
		userID int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantFile model.DataFile
		wantErr  bool
	}{
		{
			name: "find file",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				fileID: 0,
				userID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotFile, err := s.FindFile(tt.args.ctx, tt.args.fileID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFile, tt.wantFile) {
				t.Errorf("FindFile() gotFile = %v, want %v", gotFile, tt.wantFile)
			}
		})
	}
}

func TestService_FindText(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx    context.Context
		textID int
		userID int
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantText model.DataText
		wantErr  bool
	}{
		{
			name: "find text",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:    ctx,
				textID: 0,
				userID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotText, err := s.FindText(tt.args.ctx, tt.args.textID, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotText, tt.wantText) {
				t.Errorf("FindText() gotText = %v, want %v", gotText, tt.wantText)
			}
		})
	}
}

func TestService_GetRefreshToken(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "get refresh token",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:   ctx,
				token: "1111",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			_, err := s.GetRefreshToken(tt.args.ctx, tt.args.token)
			if err.Error() != "no rows in result set" {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetSignKey(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx      context.Context
		login    string
		password string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantSignKey string
		wantErr     bool
	}{
		{
			name: "get sign key",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:      ctx,
				login:    "",
				password: "",
			},
			wantSignKey: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotSignKey, err := s.GetSignKey(tt.args.ctx, tt.args.login, tt.args.password)
			if (err.Error() != "no rows in result set") != tt.wantErr {
				t.Errorf("GetSignKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSignKey != tt.wantSignKey {
				t.Errorf("GetSignKey() gotSignKey = %v, want %v", gotSignKey, tt.wantSignKey)
			}
		})
	}
}

func TestService_SaveCard(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx  context.Context
		card model.DataCard
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantId  int
		wantErr bool
	}{
		{
			name: "save card",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:  ctx,
				card: model.DataCard{},
			},
			wantErr: true,
			wantId:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotId, err := s.SaveCard(tt.args.ctx, tt.args.card)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveCard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("SaveCard() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func TestService_SaveCred(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx  context.Context
		cred model.DataCred
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantId  int
		wantErr bool
	}{
		{
			name: "save cred",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:  ctx,
				cred: model.DataCred{},
			},
			wantErr: true,
			wantId:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotId, err := s.SaveCred(tt.args.ctx, tt.args.cred)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveCred() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("SaveCred() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func TestService_SaveFile(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx  context.Context
		file model.DataFile
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantId  int
		wantErr bool
	}{
		{
			name: "save file",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:  ctx,
				file: model.DataFile{},
			},
			wantErr: true,
			wantId:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotId, err := s.SaveFile(tt.args.ctx, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("SaveFile() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func TestService_SaveText(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx  context.Context
		text model.DataText
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantId  int
		wantErr bool
	}{
		{
			name: "save text",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:  ctx,
				text: model.DataText{},
			},
			wantErr: true,
			wantId:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotId, err := s.SaveText(tt.args.ctx, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotId != tt.wantId {
				t.Errorf("SaveText() gotId = %v, want %v", gotId, tt.wantId)
			}
		})
	}
}

func TestService_SignIn(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantTokens model.Tokens
		wantErr    bool
	}{
		{
			name: "sign in",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:  ctx,
				user: model.User{Login: "1", Password: "0"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotTokens, err := s.SignIn(tt.args.ctx, tt.args.user)

			if err == storage.ErrorUserCredentials {
				assert.Error(t, err)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTokens, tt.wantTokens) {
				t.Errorf("SignIn() gotTokens = %v, want %v", gotTokens, tt.wantTokens)
			}
		})
	}
}

func TestService_SignUp(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.ReadConfig()
	if err != nil {
		t.Error(err)
	}
	tokenManager, err := auth.NewManager(cfg.JWTSecretKey)
	if err != nil {
		t.Error(err)
	}

	store := storage.New(ctx, cfg.DatabaseDsn)
	storeFiles, err := file.New(cfg.FileStorage)
	if err != nil {
		t.Error(err)
	}

	type fields struct {
		Store        storage.Interface
		StoreFiles   *file.StorageFiles
		Cfg          *config.Config
		TokenManager auth.TokenManager
	}
	type args struct {
		ctx  context.Context
		user model.User
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantTokens model.Tokens
		wantErr    bool
	}{
		{
			name: "sign up",
			fields: fields{
				Store:        store,
				StoreFiles:   storeFiles,
				Cfg:          cfg,
				TokenManager: tokenManager,
			},
			args: args{
				ctx:  ctx,
				user: model.User{Login: "1", Password: "1"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Store:        tt.fields.Store,
				StoreFiles:   tt.fields.StoreFiles,
				Cfg:          tt.fields.Cfg,
				TokenManager: tt.fields.TokenManager,
			}
			gotTokens, err := s.SignUp(tt.args.ctx, tt.args.user)

			if err == storage.ErrorUserAlreadyExists {
				assert.NoError(t, err)
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("SignUp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTokens, tt.wantTokens) {
				t.Errorf("SignUp() gotTokens = %v, want %v", gotTokens, tt.wantTokens)
			}
		})
	}
}
