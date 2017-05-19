package types

import (
	"github.com/spf13/viper"
	"os"
)

const (
	ViperValuePack   = "value-pack"
	DefaultTaxFactor = 0.65
	ValueTaxFactor   = 0.845
)

// A Product represents a resource gathered from or crafted at a Node.
type Product struct {
	Name  string  `json:"name"`
	Value uint32  `json:"value"`
	Cost  uint32  `json:"cost,omitempty"`
	Qty   float64 `json:"qty"`
}

// Profit is the expected net value of a product harvest cycle.
func (p *Product) Profit(luck float64, link bool) float64 {
	return p.ExpectedQty(luck) * p.UnitProfit(link)
}

// UnitProfit is the net value of a single product item.
func (p *Product) UnitProfit(link bool) float64 {
	var tax float64
	if viper.GetBool(ViperValuePack) {
		tax = ValueTaxFactor
	} else {
		tax = DefaultTaxFactor
	}
	var cost float64
	if link {
		cost = float64(p.Cost) * tax
	} else {
		cost = float64(p.Cost)
	}
	return float64(p.Value) * tax - cost
}

// ExpectedQty is the expected output quantity with the given luck.
func (p *Product) ExpectedQty(luck float64) float64 {
	// TODO: find actual formula - this is merely a placeholder
	return p.Qty
}
