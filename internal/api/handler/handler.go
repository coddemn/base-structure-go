package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

// Health godoc
// @Summary Проверка работоспособности сервиса
// @Description Возвращает статус OK, если сервис готов принимать запросы.
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string "Статус сервиса"
// @Router /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
