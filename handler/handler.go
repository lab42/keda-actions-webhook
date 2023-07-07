package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-github/v53/github"
	"github.com/lab42/gha-keda-webhook/counter"
	"github.com/labstack/echo/v4"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Handler struct {
	Counter counter.RedisCounter
}

func (h Handler) Probes(c echo.Context) error {
	// Return HTTP status 200 if Redis is reachable
	if h.Counter.TestConnection() {
		return c.String(http.StatusOK, "OK")
	}

	// Return HTTP status 500 if Redis is not reachable
	return c.String(http.StatusInternalServerError, "FAIL")
}

func (h Handler) Webhook(c echo.Context) error {

	// Parse the incoming event
	var payload Payload

	payloadByteArray, err := github.ValidatePayload(c.Request(), []byte(viper.GetString("SECRET_TOKEN")))
	if err != nil {
		log.Error().Msg(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	if err := json.Unmarshal(payloadByteArray, &payload); err != nil {
		log.Error().Msg(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	// Take action based on the payload.Action
	switch payload.Action {
	case "queued":
		h.Counter.Increment()
		log.Info().Interface("payload", payload).Msg("")
		return c.NoContent(http.StatusOK)
	case "in_progress":
		h.Counter.Decrement()
		log.Info().Interface("payload", payload).Msg("")
		return c.NoContent(http.StatusOK)
	default:
		return c.NoContent(http.StatusBadRequest)
	}
}
