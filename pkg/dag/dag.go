package dag

import "github.com/sirupsen/logrus"

type DAG struct {
	Nodes   []Node `json:"nodes"`
	nodeRef map[string]*Node
	entry   []string
	exit    []string
}

type Node struct {
	Name         string   `json:"name"`
	Action       string   `json:"action" default:"GET"`
	URL          string   `json:"url"`
	Dependencies []string `json:"dependencies,omitempty"`
	ParamFormat  string   `json:"paramformat,omitempty"`
}

var (
	// nolint:unused
	logger *logrus.Entry
)

func SetLogger(log *logrus.Entry) {
	logger = log
}

func NewDAG() *DAG {
	return &DAG{
		Nodes:   []Node{},
		nodeRef: make(map[string]*Node),
	}
}

func NewNode(name, action, url string, dep []string, param string) *Node {
	return &Node{
		Name:         name,
		Action:       action,
		URL:          url,
		Dependencies: dep,
		ParamFormat:  param,
	}
}
