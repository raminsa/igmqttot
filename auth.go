package mqtt

import (
	"encoding/json"
	"github.com/google/uuid"
)

type DeviceAuth struct {
	Username     string
	ClientId     string
	UserId       int64
	Password     string
	DeviceId     string
	DeviceSecret string
	Sr           string
	Rc           string
	PhoneId      string
	Json         string
}

func NewDeviceAuth(username string) *DeviceAuth {
	d := &DeviceAuth{
		Username: username,
		UserId:   0,
		PhoneId:  uuid.NewMD5(uuid.New(), []byte(username)).String(),
	}

	if len(d.PhoneId) > 20 {
		d.ClientId = d.PhoneId[:20]
	} else {
		d.ClientId = "0"
	}

	return d
}

func (d *DeviceAuth) Update() {
	if d.ClientId == "" {
		if len(d.PhoneId) > 20 {
			d.ClientId = d.PhoneId[:20]
		} else {
			d.ClientId = "0"
		}
	}
}

func (d *DeviceAuth) Read(jsonStr string) {
	if jsonStr == "" {
		return
	}

	d.Json = jsonStr

	var data map[string]any
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err == nil {
		if val, ok := data["ck"].(int64); ok {
			d.UserId = val
		}
		if val, ok := data["cs"].(string); ok {
			d.Password = val
		}
		if val, ok := data["di"].(string); ok {
			d.DeviceId = val
			if len(d.DeviceId) > 20 {
				d.ClientId = d.DeviceId[:20]
			}
		} else {
			d.DeviceId = ""
		}
		if val, ok := data["ds"].(string); ok {
			d.DeviceSecret = val
		}
		if val, ok := data["sr"].(string); ok {
			d.Sr = val
		}
		if val, ok := data["rc"].(string); ok {
			d.Rc = val
		}
	}
}

func (d *DeviceAuth) ToString() string {
	return d.Json
}
