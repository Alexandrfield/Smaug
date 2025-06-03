package worker

import (
	"encoding/json"
	"fmt"
	"time"
)

type Metadata struct {
	TypeNote   string `json:"typeNote"`
	TimeCreate string `json:"create"`
	User       string `json:"user"`
	Filename   string `json:"filename"`
}

func (met *Metadata) GetInfo() string {
	return "type:" + met.TypeNote + "; create:" + met.TimeCreate + "; user:" + met.User + ";"
}

func SerializeMetadata(dat *Metadata) []byte {
	data, _ := json.Marshal(dat)
	return data
}
func DeserializeMetadata(data []byte) (*Metadata, error) {
	var metData Metadata
	err := json.Unmarshal(data, &metData)
	if err != nil {
		return nil, fmt.Errorf("problem with deserialize Metadata from bytes. err: %w", err)
	}
	return &metData, nil
}

func NewMetadata(typeNote string, user string, create time.Time) *Metadata {
	return &Metadata{TypeNote: typeNote, User: user, TimeCreate: create.Format("2006-01-02 15:04:05")}
}
