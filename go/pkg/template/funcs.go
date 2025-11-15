package template

import (
	"encoding/json"
	"fmt"
	"strings"
)

// FormatCurrency formats a number as currency with K/M suffixes.
func FormatCurrency(value float64) string {
	if value >= 1000000 {
		return fmt.Sprintf("$%.2fM", value/1000000)
	} else if value >= 1000 {
		return fmt.Sprintf("$%.2fK", value/1000)
	}
	return fmt.Sprintf("$%.2f", value)
}

// FormatPercent formats a number as a percentage with +/- sign.
func FormatPercent(value float64) string {
	if value >= 0 {
		return fmt.Sprintf("+%.2f%%", value)
	}
	return fmt.Sprintf("%.2f%%", value)
}

// FormatFloat formats a float with specified precision.
func FormatFloat(value float64, precision int) string {
	format := fmt.Sprintf("%%.%df", precision)
	return fmt.Sprintf(format, value)
}

// ColorCode returns an emoji indicator based on sentiment.
func ColorCode(sentiment string) string {
	switch sentiment {
	case "bullish", "positive", "up":
		return "🟢"
	case "bearish", "negative", "down":
		return "🔴"
	case "neutral", "flat":
		return "🟡"
	default:
		return "⚪"
	}
}

// TrendIndicator returns an arrow indicator based on comparison.
func TrendIndicator(current, previous float64) string {
	if current > previous {
		return "📈"
	} else if current < previous {
		return "📉"
	}
	return "➡️"
}

// IsBullish checks if price is above EMA (bullish signal).
func IsBullish(price, ema float64) bool {
	return price > ema
}

// IsBearish checks if price is below EMA (bearish signal).
func IsBearish(price, ema float64) bool {
	return price < ema
}

// IsOverbought checks if RSI indicates overbought condition.
func IsOverbought(rsi float64) bool {
	return rsi > 70
}

// IsOversold checks if RSI indicates oversold condition.
func IsOversold(rsi float64) bool {
	return rsi < 30
}

// JoinFloats joins a float array into a comma-separated string.
func JoinFloats(arr []float64, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	strs := make([]string, len(arr))
	for i, v := range arr {
		strs[i] = fmt.Sprintf("%.2f", v)
	}
	return strings.Join(strs, sep)
}

// JoinInts joins an int array into a comma-separated string.
func JoinInts(arr []int, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	strs := make([]string, len(arr))
	for i, v := range arr {
		strs[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(strs, sep)
}

// JoinStrings joins a string array with a separator.
func JoinStrings(arr []string, sep string) string {
	return strings.Join(arr, sep)
}

// ToJSON converts any value to a JSON string.
func ToJSON(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(bytes)
}

// ToJSONPretty converts any value to a pretty-printed JSON string.
func ToJSONPretty(v interface{}) string {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(bytes)
}

// RangeFormat formats a range with optional unit.
func RangeFormat(min, max interface{}, unit string) string {
	return fmt.Sprintf("%v-%v%s", min, max, unit)
}

// Default returns the default value if the first value is zero/empty.
func Default(defaultValue, value interface{}) interface{} {
	// Check if value is zero/nil/empty
	switch v := value.(type) {
	case string:
		if v == "" {
			return defaultValue
		}
	case int, int64, float64:
		if v == 0 {
			return defaultValue
		}
	case bool:
		if !v {
			return defaultValue
		}
	case nil:
		return defaultValue
	}
	return value
}

// Multiply multiplies two numbers.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide divides two numbers.
func Divide(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

// Add adds two numbers.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract subtracts two numbers.
func Subtract(a, b float64) float64 {
	return a - b
}

// Abs returns the absolute value.
func Abs(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}

// Min returns the minimum of two values.
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two values.
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
