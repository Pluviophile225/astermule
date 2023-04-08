package dag

import "testing"

func TestPreflightEntryAndExit(t *testing.T) {
	A := NewNode("A", "GET", "http://10.244.1.10:8080", []string{}, "")
	B := NewNode("B", "GET", "http://10.244.1.11:8080", []string{}, "")
	C := NewNode("C", "POST", "http://10.244.1.12:8080", []string{"A", "B"}, "")
	D := NewNode("D", "GET", "http://10.244.1.13:8080", []string{"A", "C"}, "")
	E := NewNode("E", "GET", "http://10.244.1.14:8080", []string{"B"}, "")

	dag := NewDAG()
	dag.Nodes = append(dag.Nodes, *A)
	dag.Nodes = append(dag.Nodes, *B)
	dag.Nodes = append(dag.Nodes, *C)
	dag.Nodes = append(dag.Nodes, *D)
	dag.Nodes = append(dag.Nodes, *E)

	err := dag.Preflight()
	if err != nil {
		t.Error(err)
	}
	resultEntry := map[string]struct{}{"A": {}, "B": {}}
	resultExit := map[string]struct{}{"D": {}, "E": {}}
	if len(dag.entry) != 2 || len(dag.exit) != 2 {
		t.Error("Entry and Exit incorrect")
	}
	for _, nodeName := range dag.entry {
		if _, ok := resultEntry[nodeName]; !ok {
			t.Error("Entry and Exit incorrect")
		}
	}

	for _, nodeName := range dag.exit {
		if _, ok := resultExit[nodeName]; !ok {
			t.Error("Entry and Exit incorrect")
		}
	}
}
