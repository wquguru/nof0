package template

import "fmt"

// Range represents a numeric range with min and max values.
type Range struct {
	Min float64 `json:"min" doc:"Minimum value" example:"1"`
	Max float64 `json:"max" doc:"Maximum value" example:"20"`
}

// String returns a formatted string representation of the range.
func (r Range) String() string {
	return fmt.Sprintf("%.2f-%.2f", r.Min, r.Max)
}

// IsValid checks if the range is valid (max > min).
func (r Range) IsValid() bool {
	return r.Max > r.Min
}

// Contains checks if a value is within the range.
func (r Range) Contains(v float64) bool {
	return v >= r.Min && v <= r.Max
}

// Duration represents a time duration with value and unit.
type Duration struct {
	Value int    `json:"value" doc:"Duration value" example:"5"`
	Unit  string `json:"unit" doc:"Time unit (minutes, hours, days)" example:"minutes"`
}

// String returns a formatted string representation.
func (d Duration) String() string {
	return fmt.Sprintf("%d %s", d.Value, d.Unit)
}

// Minutes converts duration to minutes.
func (d Duration) Minutes() int {
	switch d.Unit {
	case "hours":
		return d.Value * 60
	case "days":
		return d.Value * 24 * 60
	case "minutes":
		return d.Value
	default:
		return d.Value
	}
}

// Percentage represents a percentage value.
type Percentage float64

// String returns a formatted percentage string.
func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", float64(p))
}

// Decimal returns the decimal representation (e.g., 5% -> 0.05).
func (p Percentage) Decimal() float64 {
	return float64(p) / 100.0
}
