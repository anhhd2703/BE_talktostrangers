package biz

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	log "talktostrangers/pkg/log"
	"talktostrangers/proto"
	"time"
)

func sendJSON(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	result := map[string]interface{}{"success": true}
	if err != nil {
		log.Infof("[sendJSON] err: %v", err)
		result["success"] = false
		result["error"] = err.Error()
		w.WriteHeader(http.StatusBadRequest)
	} else if data != nil {
		result["data"] = data
	}
	e := json.NewEncoder(w).Encode(result)
	if e != nil {
		log.Infof("[sendJSON][json.NewEncoder] err: %v", e)
	}
}

func GetRooomHandler(wh *WareHouse) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		ctx := context.Background()
		// 		//Derive a context with cancel
		ctxWithCancel, cancelFunction := context.WithTimeout(ctx, 15*time.Second)
		body := proto.UserInfo{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			log.Infof("[GetRooomHandler] [NewDecoder]", err)
			sendJSON(w, nil, err)
			return
		}
		roomId := wh.Register(body)
		if roomId != "" {
			sendJSON(w, roomId, nil)
			return
		}
		waitingRoom := wh.Sub(body.UserId)
		select {
		case roomId = <-waitingRoom:
			wh.UnRegister(body)
			wh.UnSub(body.UserId)
			sendJSON(w, roomId, nil)
			cancelFunction()
			return
		case <-ctxWithCancel.Done():
			switch ctxWithCancel.Err() {
			case context.DeadlineExceeded:
				sendJSON(w, nil, errors.New("User not found"))
				wh.UnRegister(body)
				wh.UnSub(body.UserId)
				return
			case context.Canceled:
				wh.UnRegister(body)
				wh.UnSub(body.UserId)
				log.Infof("context.Canceled")
				return
			}

		}

	})
}
