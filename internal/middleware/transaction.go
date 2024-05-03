package middleware

import (
	"1-cat-social/pkg/logger"
	"1-cat-social/pkg/response"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// StatusInList -> checks if the given status is in the list
func StatusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}

// DBTransactionMiddleware : to setup the database transaction middleware
func DBTransactionMiddleware(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		txHandle, err := db.Beginx()
		if err != nil {
			response.GenerateResponse(c, http.StatusInternalServerError, response.WithMessage(err.Error()))
			c.Abort()
			return
		}
		log.Print("beginning database transaction")

		defer func() {
			if r := recover(); r != nil {
				txHandle.Rollback()
				panic(r)
			}
		}()

		c.Set("db_trx", txHandle)
		c.Next()

		if StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			log.Print("committing transactions")
			if err := txHandle.Commit(); err != nil {
				logger.Error(err)
			}
		} else {
			logger.Info(fmt.Sprintf("rolling back transaction due to status code: %d", c.Writer.Status()))
			txHandle.Rollback()
		}
	}
}
