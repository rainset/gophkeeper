package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rainset/gophkeeper/internal/server/model"
	"github.com/rainset/gophkeeper/internal/server/service"
	"github.com/rainset/gophkeeper/internal/server/storage"
	"github.com/rainset/gophkeeper/pkg/logger"
	"net/http"
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
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.service.SignIn(rb)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, token)
}
func (h *Handler) SignUp(c *gin.Context) {

	var rb model.User
	err := c.BindJSON(&rb)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.service.SignUp(rb)
	if err != nil {
		if errors.Is(err, storage.ErrorUserAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, token)

}
func (h *Handler) RefreshToken(c *gin.Context) {

	var rb model.RefreshToken
	err := c.BindJSON(&rb)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	tokens, err := h.service.GetRefreshToken(rb.Token)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) SaveCard(c *gin.Context) {
	var err error
	var rb model.DataCard
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error(err)
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

	err = h.service.SaveCard(rb)
	if err != nil {

		if errors.Is(err, storage.ErrorRowAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		logger.Error(err, rb)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusCreated)
}
func (h *Handler) DeleteCard(c *gin.Context) {
	var err error
	var rb model.DataCard
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCard(rb.ID, userID)
	if err != nil {
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
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	card, err := h.service.FindCard(rb.ID, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, card)
}
func (h *Handler) FindAllCards(c *gin.Context) {
	var err error

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cards, err := h.service.FindAllCards(userID)
	if err != nil {
		logger.Error(err)
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
		logger.Error(err)
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

	err = h.service.SaveCred(rb)
	if err != nil {

		if errors.Is(err, storage.ErrorRowAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		logger.Error(err, rb)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusCreated)
}
func (h *Handler) DeleteCred(c *gin.Context) {
	var err error
	var rb model.DataCred
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCred(rb.ID, userID)
	if err != nil {
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
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cred, err := h.service.FindCred(rb.ID, userID)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, cred)
}
func (h *Handler) FindAllCreds(c *gin.Context) {
	var err error

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	creds, err := h.service.FindAllCreds(userID)
	if err != nil {
		logger.Error(err)
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
		logger.Error(err)
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

	err = h.service.SaveText(rb)
	if err != nil {

		if errors.Is(err, storage.ErrorRowAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		logger.Error(err, rb)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Status(http.StatusCreated)
}
func (h *Handler) DeleteText(c *gin.Context) {
	var err error
	var rb model.DataText
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.service.DeleteText(rb.ID, userID)
	if err != nil {
		logger.Error(err)
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
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	card, err := h.service.FindText(rb.ID, userID)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, card)
}
func (h *Handler) FindAllTexts(c *gin.Context) {
	var err error

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	texts, err := h.service.FindAllTexts(userID)
	if err != nil {
		logger.Error(err)
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
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	formFile, err := c.FormFile("file")
	if err != nil {
		logger.Error("file param required", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	src, err := formFile.Open()
	if err != nil {
		logger.Error("open upload file error", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	filePath, err := h.service.StoreFiles.SaveFile(src)
	if err != nil {

		if errors.Is(err, storage.ErrorRowAlreadyExists) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}

		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	file.UserID = userID
	file.Title = c.PostForm("title")
	file.Meta = c.PostForm("meta")
	file.Path = filePath
	file.Filename = formFile.Filename

	err = h.service.SaveFile(file)
	if err != nil {
		logger.Error("open upload file error", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.Status(http.StatusCreated)
}
func (h *Handler) DeleteFile(c *gin.Context) {
	var err error
	var rb model.DataFile
	err = c.BindJSON(&rb)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = h.service.DeleteFile(rb.ID, userID)
	if err != nil {
		logger.Error(err)
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
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	file, err := h.service.FindFile(rb.ID, userID)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, file)
}
func (h *Handler) FindAllFiles(c *gin.Context) {
	var err error

	userID, err := h.getUserIDFromRequest(c)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	files, err := h.service.FindAllFiles(userID)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if len(files) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	c.JSON(http.StatusOK, files)
}
