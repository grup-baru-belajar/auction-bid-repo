package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/services"
	"github.com/grup-baru-belajar/auction-bid-repo/internal/token"
)

func respondBindError(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		missing := make([]string, 0, len(ve))
		for _, fe := range ve {
			missing = append(missing, strings.ToLower(fe.Field()))
		}
		respondError(c, http.StatusBadRequest,
			fmt.Sprintf("Missing or invalid field: %s", strings.Join(missing, ", ")))
		return
	}

	var se *json.SyntaxError
	var ute *json.UnmarshalTypeError
	if errors.As(err, &se) || errors.As(err, &ute) {
		respondError(c, http.StatusBadRequest, "Malformed JSON body")
		return
	}

	respondError(c, http.StatusBadRequest, "Invalid request body")
}

func respondWithError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidCredentials):
		respondError(c, http.StatusUnauthorized, "Invalid username or password")

	case errors.Is(err, token.ErrInvalidToken):
		respondError(c, http.StatusUnauthorized, "Invalid or expired token")

	case errors.Is(err, services.ErrAuctionNotFound):
		respondError(c, http.StatusNotFound, "Auction not found")

	case errors.Is(err, services.ErrInvalidStartingPrice):
		respondError(c, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrInvalidEndTime):
		respondError(c, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrBidTooLow):
		respondError(c, http.StatusBadRequest, err.Error())

	case errors.Is(err, services.ErrAuctionCompleted):
		respondError(c, http.StatusConflict, "Auction already completed")

	default:
		log.Printf("unhandled error on %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
		respondError(c, http.StatusInternalServerError, "Internal server error")
	}

}
