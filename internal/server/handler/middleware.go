package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rainset/gophkeeper/pkg/logger"
)

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.service.TokenManager.Parse(headerParts[1])
}

func (h *Handler) authMiddleware(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		logger.Info("authMiddleware:", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("user_id", id)
}

func (h *Handler) getUserIDFromRequest(c *gin.Context) (userID int, err error) {
	ctxUserID, ex := c.Get("user_id")

	if !ex {
		return userID, errors.New("failed to convert `user_id` from ctx")
	}

	userID, err = strconv.Atoi(ctxUserID.(string))
	if err != nil {
		return userID, errors.New("failed to parse `user_id` from ctx")
	}

	return userID, err
}
