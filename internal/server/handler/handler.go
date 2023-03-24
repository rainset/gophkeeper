package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/internal/server/service"
	"github.com/rainset/gophkeeper/internal/server/storage"
	"github.com/rainset/gophkeeper/pkg/logger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Init() *gin.Engine {
	r := gin.Default()

	r.StaticFS("/_file_storage", gin.Dir(h.service.Cfg.FileStorage, false))

	r.MaxMultipartMemory = 16 << 20 // 16 MiB

	r.GET("/ping", h.Ping)
	r.POST("/sign-up", h.SignUp)
	r.POST("/sign-in", h.SignIn)

	r.POST("/refresh-token", h.RefreshToken)
	r.POST("/sign-key", h.SignKey)

	store := r.Group("/store", h.authMiddleware)
	{
		store.POST("/card", h.SaveCard)
		store.DELETE("/card", h.DeleteCard)
		store.GET("/card", h.FindCard)
		store.GET("/card/list", h.FindAllCards)

		store.POST("/cred", h.SaveCred)
		store.DELETE("/cred", h.DeleteCred)
		store.GET("/cred", h.FindCred)
		store.GET("/cred/list", h.FindAllCreds)

		store.POST("/text", h.SaveText)
		store.DELETE("/text", h.DeleteText)
		store.GET("/text", h.FindText)
		store.GET("/text/list", h.FindAllTexts)

		store.POST("/file", h.SaveFile)
		store.DELETE("/file", h.DeleteFile)
		store.GET("/file", h.FindFile)
		store.GET("/file/list", h.FindAllFiles)
	}

	return r
}

func (h *Handler) Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

func (h *Handler) SignIn(c *gin.Context) {
	var rb model.User
	err := c.BindJSON(&rb)

	if err != nil {
		logger.Error("SignIn Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	token, err := h.service.SignIn(c, rb)
	if err != nil {
		logger.Error("SignIn Handler: ", err, rb)
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *Handler) SignUp(c *gin.Context) {
	var rb model.User
	err := c.BindJSON(&rb)
	if err != nil {
		logger.Error("SignUp Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	token, err := h.service.SignUp(c, rb)
	if err != nil {
		if errors.Is(err, storage.ErrorUserAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)

			return
		}

		logger.Error("SignUp Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var rb model.RefreshToken
	err := c.BindJSON(&rb)
	if err != nil {
		logger.Error("RefreshToken Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	tokens, err := h.service.GetRefreshToken(c, rb.Token)
	if err != nil {
		logger.Error("RefreshToken Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) SignKey(c *gin.Context) {
	var rb model.User
	err := c.BindJSON(&rb)

	if err != nil {
		logger.Error("SignKey Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	signKey, err := h.service.GetSignKey(c, rb.Login, rb.Password)
	if err != nil {
		logger.Error("SignKey Handler: ", err, rb)
		c.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	c.JSON(http.StatusOK, gin.H{"sign_key": signKey})
}

func (h *Handler) SaveCard(c *gin.Context) {
	var err error
	var rb model.DataCard

	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("SaveCard Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("SaveCard Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	rb.UserID = userID

	id, err := h.service.SaveCard(c, rb)
	if err != nil {
		if errors.Is(err, storage.ErrorRowAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)

			return
		}

		logger.Error("SaveCard Handler: ", err, rb)

		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) DeleteCard(c *gin.Context) {
	var err error
	var rb model.DataCard

	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("DeleteCard Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("DeleteCard Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	err = h.service.DeleteCard(c, rb.ID, userID)
	if err != nil {
		logger.Error("DeleteCard Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) FindCard(c *gin.Context) {
	var err error
	var rb model.DataCard
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("FindCard Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	card, err := h.service.FindCard(c, rb.ID, userID)
	if err != nil {
		logger.Error("FindCard Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *Handler) FindAllCards(c *gin.Context) {
	var err error

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("FindAllCards Handler: ", err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	cards, err := h.service.FindAllCards(c, userID)
	if err != nil {
		logger.Error("FindAllCards Handler: ", err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	if len(cards) == 0 {
		c.Status(http.StatusNoContent)

		return
	}

	c.JSON(http.StatusOK, cards)
}

func (h *Handler) SaveCred(c *gin.Context) {
	var err error
	var rb model.DataCred
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("SaveCred Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	rb.UserID = userID

	id, err := h.service.SaveCred(c, rb)
	if err != nil {
		if errors.Is(err, storage.ErrorRowAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)

			return
		}

		logger.Error("SaveCred Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})

}

func (h *Handler) DeleteCred(c *gin.Context) {
	var err error
	var rb model.DataCred
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("DeleteCred Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("DeleteCred Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	err = h.service.DeleteCred(c, rb.ID, userID)
	if err != nil {
		logger.Error("DeleteCred Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) FindCred(c *gin.Context) {
	var err error
	var rb model.DataCred
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("FindCred Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("FindCred Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	cred, err := h.service.FindCred(c, rb.ID, userID)
	if err != nil {
		logger.Error("FindCred Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, cred)
}

func (h *Handler) FindAllCreds(c *gin.Context) {
	var err error

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("FindAllCreds Handler: ", err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	creds, err := h.service.FindAllCreds(c, userID)
	if err != nil {
		logger.Error("FindAllCreds Handler: ", err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	if len(creds) == 0 {
		c.Status(http.StatusNoContent)

		return
	}

	c.JSON(http.StatusOK, creds)
}

func (h *Handler) SaveText(c *gin.Context) {
	var err error
	var rb model.DataText
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("SaveText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("SaveText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	rb.UserID = userID

	id, err := h.service.SaveText(c, rb)
	if err != nil {
		if errors.Is(err, storage.ErrorRowAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)

			return
		}

		logger.Error("SaveText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *Handler) DeleteText(c *gin.Context) {
	var err error
	var rb model.DataText
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("DeleteText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("DeleteText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	err = h.service.DeleteText(c, rb.ID, userID)
	if err != nil {
		logger.Error("DeleteText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) FindText(c *gin.Context) {
	var err error
	var rb model.DataText
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("FindText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("FindText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	card, err := h.service.FindText(c, rb.ID, userID)
	if err != nil {
		logger.Error("FindText Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, card)
}

func (h *Handler) FindAllTexts(c *gin.Context) {
	var err error

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("FindAllTexts Handler: ", err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	texts, err := h.service.FindAllTexts(c, userID)
	if err != nil {
		logger.Error("FindAllTexts Handler: ", err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	if len(texts) == 0 {
		c.Status(http.StatusNoContent)

		return
	}

	c.JSON(http.StatusOK, texts)
}

func (h *Handler) SaveFile(c *gin.Context) {
	var file model.DataFile

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("SaveFile Handler: ", err, file)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	formFile, err := c.FormFile("file")
	if err != nil {
		logger.Error("SaveFile Handler: ", err, file)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	src, err := formFile.Open()
	if err != nil {
		logger.Error("SaveFile Handler: ", err, file)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	filePath, err := h.service.StoreFiles.SaveFile(src)
	if err != nil {
		if errors.Is(err, storage.ErrorRowAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)

			return
		}

		logger.Error("SaveFile Handler: ", err, file)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	t, err := time.Parse(time.RFC3339, c.PostForm("updated_at"))
	if err != nil {
		logger.Error("SaveFile Handler updated_at format RFC3339 error: ", err, file)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		logger.Error("SaveFile Handler parse id error: ", err, file)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if id > 0 {
		file.ID = id
	}

	file.UserID = userID
	file.Title = c.PostForm("title")
	file.Meta = c.PostForm("meta")
	file.Path = filePath
	file.UpdatedAt = t
	file.Filename = formFile.Filename

	fileID, err := h.service.SaveFile(c, file)
	if err != nil {
		logger.Error("SaveFile Handler open upload file error: ", err, file)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": fileID})
}

func (h *Handler) DeleteFile(c *gin.Context) {
	var err error
	var rb model.DataFile
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error("DeleteFile Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("DeleteFile Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	err = h.service.DeleteFile(c, rb.ID, userID)
	if err != nil {
		logger.Error("DeleteFile Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) FindFile(c *gin.Context) {
	var err error
	var rb model.DataFile
	err = c.BindJSON(&rb)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("FindFile Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	file, err := h.service.FindFile(c, rb.ID, userID)
	if err != nil {
		logger.Error("FindFile Handler: ", err, rb)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	c.JSON(http.StatusOK, file)
}

func (h *Handler) FindAllFiles(c *gin.Context) {
	var err error

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error("FindAllFiles Handler: ", err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	files, err := h.service.FindAllFiles(c, userID)
	if err != nil {
		logger.Error("FindAllFiles Handler: ", err)
		c.AbortWithStatus(http.StatusBadRequest)

		return
	}

	if len(files) == 0 {
		c.Status(http.StatusNoContent)

		return
	}

	c.JSON(http.StatusOK, files)
}
