package parser

import (
	// "encoding/json"

	"github.com/Pluviophile225/astermule/pkg/dag"
	"github.com/sirupsen/logrus"
)

var (
	logger *logrus.Entry
)

func SetLogger(log *logrus.Entry) {
	logger = log
}

type ControlPlane struct {
	Fs    []func()
	Entry []chan<- Message
	Exit  []<-chan Message
}

type Parser interface {
	Parse(*dag.DAG) ControlPlane
}

type Message struct {
	Status Status `json:"status"`
}

// TODO: Define Status
type Status struct {
	Health bool `json:"health"`
}

func NewMessage(health bool, data string) *Message {
	return &Message{
		Status: Status{
			Health: health,
		},
	}
}

// func (m Message) Marshal() ([]byte, error) {
// 	return json.Marshal(m)
// }

// func (m Message) Unmarshal() (interface{}, error) {
// 	var data interface{}
// 	err := json.Unmarshal([]byte(m.Data), &data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
