package patient

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"patients/pkg/logging"
	"time"
)

const contextTimeOut = time.Second * 10

type Handler struct {
	logger *logging.Logger
	repo   Storage
}

func NewHandler(repo Storage, logger *logging.Logger) *Handler {
	return &Handler{
		logger: logger,
		repo:   repo,
	}
}

func (h *Handler) NewPatient(c *gin.Context) {
	var person Person

	if err := c.ShouldBindJSON(&person); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	defer cancel()

	err := h.repo.Create(ctx, &person)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
		default:
			h.logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, person)
}

func (h *Handler) GetListPatients(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	defer cancel()

	persons, err := h.repo.FindAll(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
		default:
			h.logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, persons)
}

func (h *Handler) EditPatient(c *gin.Context) {

	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	defer cancel()

	err := h.repo.Update(ctx, &person)
	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
		default:
			h.logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, person)
}

func (h *Handler) DelPatient(c *gin.Context) {
	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), contextTimeOut)
	defer cancel()

	err := h.repo.Delete(ctx, person.Guid)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			c.JSON(http.StatusNotFound, gin.H{"error": "patient not found"})
		case errors.Is(err, context.DeadlineExceeded):
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
		default:
			h.logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, "deleted")
}
