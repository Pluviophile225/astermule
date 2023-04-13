package parser

import (
	"encoding/json"

	"github.com/kasterism/astermule/pkg/clients/httpclient"
	"github.com/kasterism/astermule/pkg/dag"

	"github.com/kasterism/astermule/tools/whiteboard"
)

type SimpleParser struct {
	ChanGroup map[string]*ChannelGroup
	ActionMap map[string]map[string]string
}

type ChannelGroup struct {
	ReadCh  []<-chan Message
	WriteCh []chan<- Message
}

func NewSimpleParser(action map[string]map[string]string) *SimpleParser {
	return &SimpleParser{
		ChanGroup: make(map[string]*ChannelGroup),
		ActionMap: action,
	}
}

func (s *SimpleParser) Parse(d *dag.DAG) (ControlPlane, []string) {
	c := ControlPlane{}
	s.Init(d)
	exit := s.makeChannelGroup(d)
	c.Entry, c.Exit = s.scanChannelGroup(d)
	c.Fs = s.makeFunc(d)
	return c, exit
}

func (s *SimpleParser) Init(d *dag.DAG) {
	for _, node := range d.Nodes {
		s.ChanGroup[node.Name] = &ChannelGroup{
			ReadCh:  make([]<-chan Message, 0),
			WriteCh: make([]chan<- Message, 0),
		}
	}
}

func (s *SimpleParser) makeChannelGroup(d *dag.DAG) []string {
	for _, node := range d.Nodes {
		for _, dep := range node.Dependencies {
			ch := make(chan Message)
			s.ChanGroup[dep].WriteCh = append(s.ChanGroup[dep].WriteCh, ch)
			s.ChanGroup[node.Name].ReadCh = append(s.ChanGroup[node.Name].ReadCh, ch)
		}
	}
	var exit []string
	for _, node := range d.Nodes {
		if len(s.ChanGroup[node.Name].WriteCh) == 0 {
			exit = append(exit, node.Name)
		}
	}
	return exit

}

func (s *SimpleParser) scanChannelGroup(d *dag.DAG) ([]chan<- Message, []<-chan Message) {
	entry := make([]chan<- Message, 0)
	exit := make([]<-chan Message, 0)
	for _, v := range s.ChanGroup {
		if len(v.ReadCh) == 0 {
			ch := make(chan Message)
			entry = append(entry, ch)
			v.ReadCh = append(v.ReadCh, ch)
		}

		if len(v.WriteCh) == 0 {
			ch := make(chan Message)
			exit = append(exit, ch)
			v.WriteCh = append(v.WriteCh, ch)
		}
	}
	return entry, exit
}

func (s *SimpleParser) makeFunc(d *dag.DAG) []func() {
	fs := make([]func(), 0)
	for i := range d.Nodes {
		node := d.Nodes[i]
		chGrp := s.ChanGroup[node.Name]
		f := func() {
			for {
				logger.Infoln("func register:", node.Name)
				for _, readCh := range chGrp.ReadCh {
					<-readCh
				}

				logger.Infoln("func launch:", node.Name)

				// TODO
				// This is where the method is called to get the parameters from the ActionMap
				// data = fun(ActionMap,node.ParamFormat)
				data, err := whiteboard.Bend(s.ActionMap, node.ParamFormat, nil)
				_data := data.(string)
				// Prepare sendMsg
				sendMsg := &Message{}
				sendMsg.Status.Health = true

				// Call http client
				logger.Infoln("send msg to", node.URL)

				// Determine whether to use a Get or Post request
				res, err := httpclient.Send(node.Action, node.URL, _data)
				if err != nil {
					logger.Errorln("httpclient error:", err)
					sendMsg.Status.Health = false
				} else {
					logger.Infoln("receive respense:", res)
					var result map[string]string
					err := json.Unmarshal([]byte(res), &result)
					if err != nil {
						logger.Fatalln("Entry Param wrong, json transform:", err)
					}
					s.ActionMap[node.Name] = result
				}

				for _, writeCh := range chGrp.WriteCh {
					writeCh <- *sendMsg
				}

				logger.Infoln("func end:", node.Name)
			}
		}
		fs = append(fs, f)
	}
	return fs
}
