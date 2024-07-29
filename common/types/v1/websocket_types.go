package v1

import (
	"encoding/json"
	"sync"

	"golang.org/x/net/websocket"
	"k8s.io/klog/v2"
)

type eventAlarmMsg struct {
	CType  string `json:"type"`
	Detail string `json:"detail"`
	Level  int64  `json:"level"`
}

var (
	Connections = make(map[*websocket.Conn]bool)
	WS_Mutex    sync.Mutex
)

func SendEventMsgToWeb(ctype, detail string, level int64) error {
	WS_Mutex.Lock()
	defer WS_Mutex.Unlock()
	msg := &eventAlarmMsg{
		CType:  ctype,
		Detail: detail,
		Level:  level,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		klog.Errorf("json marshal error: %s", err)
		return err
	}
	for conn := range Connections {
		err := websocket.Message.Send(conn, string(data))
		if err != nil {
			klog.Errorf("send msg by websocket error: %v", err.Error())
			delete(Connections, conn)
			continue
		}
	}
	return nil
}
