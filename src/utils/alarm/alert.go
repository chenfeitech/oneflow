package alarm

import (
	"encoding/json"
	log "github.com/cihub/seelog"
)

type AlarmMessage struct {
	Catalog   string                  `json:"catalog"`
	Title     string                  `json:"title"`
	Message   string                  `json:"message"`
	ExtraData *map[string]interface{} `json:"extra_data"`
}

func RaiseAlarm(catalog string, title string, message string, extra_data *map[string]interface{}) {
	var json_bytes []byte
	var err error
	alarm_msg := AlarmMessage{catalog, title, message, extra_data}
	if json_bytes, err = json.Marshal(alarm_msg); err != nil {
		log.Error("Marshal AlarmMessage failed:", err)
		alarm_msg.ExtraData = nil
		json_bytes, err = json.Marshal(alarm_msg)
	}
	if err == nil {
		log.Error("IDATA_ALARM:", string(json_bytes))
	} else {
		log.Error("Marshal AlarmMessage failed:", err)
		log.Error("IDATA_ALARM:", catalog, " title:", title, " message:", message, " extra data:", extra_data)
	}
}
