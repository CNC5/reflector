package log

import (
	"encoding/json"
	"fmt"
	"time"
)

type MsgChain struct {
	messages map[string]any
	logLevel int
	msgLevel int
	callback func(map[string]any)
}

func (mc *MsgChain) Update(key string, value any) *MsgChain {
	mc.messages[key] = value
	return mc
}

func (mc *MsgChain) UpdateWithJSON(jsonString string) *MsgChain {
	jsonMap := &map[string]string{}
	reservedKeys := map[string]bool{
		"level": true,
	}
	json.Unmarshal([]byte(jsonString), jsonMap)
	for key, value := range *jsonMap {
		if _, exists := reservedKeys[key]; exists {
			continue
		}
		mc.messages[key] = value
	}
	return mc
}

func (mc *MsgChain) Done() {
	if !(mc.logLevel >= mc.msgLevel) {
		return
	}
	mc.Update("timestamp", time.Now().UTC().Format(time.UnixDate))
	mc.callback(mc.messages)
}

func (mc *MsgChain) Msg(msg string) {
	mc.messages["msg"] = msg
	mc.Done()
}

func (mc *MsgChain) Msgf(format string, a ...any) {
	mc.messages["msg"] = fmt.Sprintf(format, a...)
	mc.Done()
}
