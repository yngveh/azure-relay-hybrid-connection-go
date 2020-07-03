package arhc

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func (c *Client) Listen(handler func(*Response, *Request) error) error {

	params := url.Values{}
	params.Add("sb-hc-action", "listen")
	params.Add("sb-hc-token", c.token.Token)

	u := url.URL{
		Scheme:   "wss",
		Host:     fmt.Sprintf("%s:443", c.config.Namespace),
		Path:     fmt.Sprintf("$hc/%s", c.config.ConnectionName),
		RawQuery: params.Encode(),
	}

	debug("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	done := make(chan struct{})
	go readLoop(conn, handler, done)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:

			debug("Interrupted send close message")

			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Print("write close:", err)
				return err
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}

			return nil
		}
	}
}

func readLoop(conn *websocket.Conn, handler func(*Response, *Request) error, done chan struct{}) {

	defer close(done)

	for {

		debug("Waiting for message")

		t, message, err := conn.ReadMessage()
		if err != nil {
			debug("read:", err)
			break
		}

		debug("MessageType: %d, Message: %s", t, string(message))

		requestObj := &RequestObj{}
		if err := json.Unmarshal(message, requestObj); err != nil {
			debug("error unmarshal:", err)
			continue
		}

		request := requestObj.Request

		debug("Request target: %s", request.Target)

		if request.Body {
			debug("Reading body frame")
			t, body, err := conn.ReadMessage()
			if err != nil {
				debug("read:", err)
				continue
			}
			debug("Type: %s, Body: %s", t, string(body))
			request.SetRequestBody(body)
		}

		response := Response{
			RequestID:  request.ID,
			StatusCode: "200",
		}

		debug("Invoking handler")
		err = handler(&response, &request)
		if err != nil {
			return
		}

		responseObj := &ResponseObj{
			Response: response,
		}

		debug("Writing header frame")
		if err := conn.WriteJSON(responseObj); err != nil {
			debug("error unmarshal:", err)
			continue
		}

		if response.Body {
			debug("Writing body frame")
			if err := conn.WriteMessage(websocket.BinaryMessage, response.responseBody); err != nil {
				debug("error writing body:", err)
			}
		}
	}
}
