package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
)

const ViperNodeMap = "node-map"

func init() {
	viper.SetDefault(ViperNodeMap, make(NodeMap, 0))
}

// A Node represents a location on the map.
type Node struct {
	Name     string         `json:"name,omitempty"`
	Cost     int            `json:"cost"`
	Distance map[string]int `json:"distance"`

	A []string           `json:"adjacent"`
	P map[string]float64 `json:"products,omitempty"`

	Adjacent []*Node    `json:"-"`
	Products []*Product `json:"-"`

	Worker *Worker `json:"-"`
}

func ProfitPerHour(p *Product, w *Worker, work float64, distance uint16) float64 {
	return p.Profit(w.Luck, false)/w.Time(work, distance)
}

// Root is the name of a root node
type Root string
const (
	Velia  Root = "velia"
	Olvia  Root = "olvia"
	Hiedel Root = "hiedel"
	//...
)

type Node interface {
	// Name of the node
	Name() string
	// Cost to connect the node
	Cost() int
	// Prev node in the path to root
	Prev(root Root) Node

	priority() int
	setPriority(p int)
	index() int
	setIndex(i int)
}
type node struct {
	N string        `json:"name"`
	C int           `json:"cost"`
	P map[Root]Node `json:"prev"`

	p int
	i int
}
func (n node) Name() string        { return n.N }
func (n node) Cost() int           { return n.C }
func (n node) Prev(root Root) Node { return n.P[root] }

type NodeHeap []*Node
func (h NodeHeap) Len() int { return len(h) }
func (h NodeHeap) Less(i, j int) bool {
	return h[i].priority() > h[j].priority()
}
func (h NodeHeap) Swap(i, j int) bool {
	h[i], h[j] = h[j], h[i]
	h[i].setIndex(i)
	h[j].setIndex(j)
}
func (h *NodeHeap) Push(x interface{}) {
	item := Node(x)
	item.setIndex(len(*h))
	*h = append(*h, item)
}
func (h *NodeHeap) Pop() interface{} {
	item := (*h)[len(*h)-1]
	item.setIndex(-1)
	*h = (*h)[0:len(*h)-1]
	return item
}

// Path connecting the given Node and Root.
func Path(n Node, root Root) []Node {
	nodes := make([]Node, 0)
	prev := n.Prev(root)
	for prev != nil {
		nodes = append(nodes, prev)
		prev = prev.Prev(root)
	}
	return nodes
}

// PathCost for connecting the given Node and Root.
func PathCost(n Node, root Root) int {
	var cost int
	for _, p := range Path(n, root) {
		cost += p.Cost()
	}
	return cost
}

// A RootNode is a starting point for node connections and workers.
type RootNode struct {
	Name       string
	WorkerCost []int
	Workers    []*Worker
}

// A ConnNode connects nodes together, but has no output of its own.
type ConnNode struct {
	Node
}

// A WorkNode can be used by workers to harvest or craft products.
type WorkNode struct {
	Node
	Dist     map[string]uint16  // Distance from each source
	Products []*Product
}

func (n *Node) Cost(source string) uint16 {
	return n.costFunc(source, 0)
}
func (n *Node) costFunc(source string, cost uint16) uint16 {
	prev, ok := n.Prev[source]
	if !ok {
		return cost
	}
	// Add this node's CP to the total cost and follow connection
	return prev.costFunc(source, cost + n.CP)
}

/* First worker at each home is free
 * For each home, each new worker costs a given CP
 * Each home has a max of N workers
 * 
 * Connecting a node to a home costs X cp
 * If multiple nodes share a partial path, they share that CP cost
 * 	(1) Split evenly by all nodes
 * 	(2) Cost is associated with the highest value node
 * 	(3) Cost is weighted as a combination of the above
 * 
 * [*] -- [w]
 *     \- [x] -- [y]
 * 			  \- [z]
 * wCP = [0, 2, 3, 1, 2, 3, 1]
 *
 * Cost of tree is w.CP + x.CP + y.CP + z.CP + SUM(wCP[*][:3])
 * Cost of adding node v = v.CP + wCP[4]
 * Cost of adding node u = u.CP + wCP[4]
 * 	Once one of the above is added, cost of node t = t.CP + wCP[5]
 * 	etc.
 * 
 * relative value of connecting a given node N to H: 
 * 	cpN = N.CP + N.Prev.CP + N.Prev.Prev.CP ... until connected to H
 *  cpW = WorkerCP[H][n] where n is the new number of connected nodes
 * 	Value = N.Profit[H]/(cpN + cpW)
 * If Value > branchValue, branchValue += Value
 */

// A NodeMap represents a graph of connected Nodes.
type NodeMap map[string]*Node

func ImportNodeMap(filename string) error {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	var nodeMap NodeMap
	err = yaml.Unmarshal(raw, nodeMap)
	if err != nil {
		return err
	}
	viper.Set(ViperNodeMap, nodeMap)
	return nil
}

func GetNodeMap() NodeMap {
	return NodeMap(viper.Get(ViperNodeMap))
}

// MarshalJSON serializes the NodeMap into raw JSON.
func (nodeMap NodeMap) MarshalJSON() ([]byte, error) {
	// Convert the NodeMap into a raw map
	raw := map[string]*Node{}
	for name, node := range nodeMap {
		// Strip the node name
		n := *node
		n.Name = ""
		raw[name] = &n
	}
	// Marshal the raw map data
	return json.Marshal(raw)
}

// UnmarshalJSON parses the given raw JSON into an initialized NodeMap.
func (nodeMap NodeMap) UnmarshalJSON(b []byte) (err error) {
	// Unmarshal the raw map data
	raw := map[string]*Node{}
	err = json.Unmarshal(b, &raw)
	if err != nil {
		return
	}
	// Convert the raw map into a NodeMap
	for name, node := range raw {
		// Trim and update the node name
		tname := trimName(name)
		if node.Name == "" {
			node.Name = tname
		}
		nodeMap[tname] = node
	}
	for name, node := range nodeMap {
		// Populate the adjacent node pointers from the map
		node.Adjacent = make([]*Node, len(node.A))
		for i, adj := range node.A {
			adj = trimName(adj)
			n, ok := nodeMap[adj]
			if !ok {
				err = fmt.Errorf("%s: unrecognized node %s", name, adj)
				return
			}
			node.Adjacent[i] = n
		}
	}
	return
}

func trimName(name string) string {
	replacer := strings.NewReplacer(" ", "_")
	return replacer.Replace(strings.ToLower(strings.TrimSpace(name)))
}
