package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/v53/github"
	"github.com/lab42/gha-keda-webhook/counter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
	var payload map[string]interface{}

	payloadByteArray, err := github.ValidatePayload(c.Request(), []byte(viper.GetString("SECRET_TOKEN")))
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := json.Unmarshal(payloadByteArray, &payload); err != nil {
		log.Infof(err.Error())
		return c.NoContent(http.StatusBadRequest)
	}

	// Take action based on the payload.Action
	p, _ := json.Marshal(payload)
	switch payload["action"] {
	case "queued":
		h.Counter.Increment()
		log.Infof("Webhook event processed successfully", string(p))
		log.Infof("payload: %s", string(p))
		return c.NoContent(http.StatusOK)
	case "in_progress":
		h.Counter.Decrement()
		log.Infof("Webhook event processed successfully", string(p))
		log.Infof("payload: %s", string(p))
		return c.NoContent(http.StatusOK)
	default:
		return c.NoContent(http.StatusBadRequest)
	}
}

func PrettyPrint(data interface{}) {
	var p []byte
	//    var err := error
	p, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s \n", p)
}
