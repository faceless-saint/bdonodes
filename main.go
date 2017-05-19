package main

import (
	"github.com/faceless-saint/bdonodes/types"
	"github.com/ghodss/yaml"
	"fmt"
)

func main() {
	n := types.NodeMap{
		"node_name": &types.Node{
			Distance: map[string]uint16{
				"velia": 1000,
				"olvia": 600,
			},
			Cost: 2,
			A: []string{"link_node"},
			P: []string{"iron_ore"},
		},
		"link_node": &types.Node{
			Distance: map[string]uint16{
				"velia": 1300,
				"olvia": 200,
			},
			Cost: 1,
			A: []string{"node_name", "olvia"},
		},
		"olvia": &types.Node{
			Distance: map[string]uint16{
				"velia": 1500,
				"olvia": 0,
			},
			Cost: 0,
			A: []string{"link_node"},
		},
	}
	b, err := yaml.Marshal(n)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(b))
	err = yaml.Unmarshal(b, &n)
	if err != nil {
		panic(err)
	}
	for name, node := range n {
		fmt.Printf("%s: %v (%d)\n", name, node.Distance, node.Cost)
		for _, n := range node.Adjacent {
			fmt.Printf("\t(%d) > %s\n", n.Cost, n.Name)
		}
	}
}
