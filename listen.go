package arhc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

//func (c *Client) HandleFunc(handler func(*ResponseWriter, *Request) error) error {
//
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt)
//
//	params := url.Values{}
//	params.Add("sb-hc-action", "listen")
//	params.Add("sb-hc-token", c.token.Token)
//
//	u := url.URL{
//		Scheme:   "wss",
//		Host:     fmt.Sprintf("%s:443", c.config.Namespace),
//		Path:     fmt.Sprintf("$hc/%s", c.config.ConnectionName),
//		RawQuery: params.Encode(),
//	}
//
//	debug("connecting to %s", u.String())
//
//	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
//	if err != nil {
//		return err
//	}
//	defer func() {
//		_ = conn.Close()
//	}()
//
//	done := make(chan struct{})
//
//	go func() {
//		defer close(done)
//		for {
//
//			t, message, err := conn.ReadMessage()
//			if err != nil {
//				log.Println("read:", err)
//				return
//			}
//
//			debug("MessageType: %d, Message: %s", t, string(message))
//
//			request := &RequestObj{}
//			if err := json.Unmarshal(message, request); err != nil {
//				log.Println("error unmarshal:", err)
//				return
//			}
//
//			debug("request: %s", request.Request.Target)
//
//			if request.Request.Body {
//				_, body, err := conn.ReadMessage()
//				log.Print("BODY:", string(body))
//				if err != nil {
//					log.Println("read:", err)
//					return
//				}
//			}
//
//			response := Response{
//				RequestID:  request.Request.ID,
//				StatusCode: "200",
//				Body:       true,
//			}
//
//			responseObj := &ResponseObj{
//				Response: response,
//			}
//
//			rw := &ResponseWriter{}
//			if err := handler(rw, &request.Request); err != nil {
//				log.Println("error unmarshal:", err)
//				return
//			}
//
//			if err := conn.WriteJSON(responseObj); err != nil {
//				log.Println("error unmarshal:", err)
//				return
//			}
//
//			if err := conn.WriteMessage(websocket.BinaryMessage, rw.b); err != nil {
//				log.Println("error writing body:", err)
//			}
//		}
//	}()
//
//	for {
//		select {
//		case <-done:
//			return nil
//		case <-interrupt:
//			log.Println("interrupt")
//
//			Cleanly close the connection by sending a close message and then
//			waiting (with timeout) for the server to close the connection.
//			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
//			if err != nil {
//				log.Println("write close:", err)
//				return
//			}
//select {
//case <-done:
//case <-time.After(time.Second):
//}
//return nil
//}
//}

//}

func (c *Client) Listen(ctx context.Context, handler func(*Response, *Request) error) error {

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

	go readLoop(conn, handler)
	go shutdown(ctx, conn)

	return nil
}

func shutdown(ctx context.Context, conn *websocket.Conn) {
	debug("Waiting for done")

	<-ctx.Done()

	debug("Closing ctx")

	err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}

	debug("Closing connection")
	_ = conn.Close()
}

func readLoop(conn *websocket.Conn, handler func(*Response, *Request) error) {

	defer func() {
		_ = conn.Close()
	}()

	for {
		t, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		debug("MessageType: %d, Message: %s", t, string(message))

		request := &RequestObj{}
		if err := json.Unmarshal(message, request); err != nil {
			log.Println("error unmarshal:", err)
			return
		}

		debug("request: %s", request.Request.Target)

		if request.Request.Body {
			_, body, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			debug("Body: %s", string(body))
			request.Request.SetRequestBody(body)
		}

		response := Response{
			RequestID:  request.Request.ID,
			StatusCode: "200",
		}

		err = handler(&response, &request.Request)
		if err != nil {
			return
		}

		responseObj := &ResponseObj{
			Response: response,
		}

		if err := conn.WriteJSON(responseObj); err != nil {
			log.Println("error unmarshal:", err)
			return
		}

		if response.Body {
			if err := conn.WriteMessage(websocket.BinaryMessage, response.responseBody); err != nil {
				log.Println("error writing body:", err)
			}
		}
	}
}
