package template_test

import (
	"testing"

	"nof0-api/pkg/template"
)

// TestSystemPromptRender tests rendering the system prompt template.
func TestSystemPromptRender(t *testing.T) {
	// Create sample data
	data := template.SystemPromptData{
		Model: template.ModelConfig{
			Name: "GPT-4-Turbo",
		},
		Market: template.MarketConfig{
			Exchange:        "Hyperliquid",
			AssetUniverse:   "BTC, ETH, SOL, AVAX",
			StartingCapital: 10000.00,
			MarketHours:     "24/7",
			ContractType:    "Perpetual futures",
			Leverage: template.Range{
				Min: 1,
				Max: 20,
			},
			TradingFee: template.Range{
				Min: 0.02,
				Max: 0.05,
			},
			Slippage: template.Range{
				Min: 0.1,
				Max: 0.5,
			},
			MinPositionSize:          100.00,
			MaxPositionConcentration: template.Percentage(30),
		},
		Risk: template.RiskConfig{
			MaxLossPerTrade: template.Range{
				Min: 1,
				Max: 3,
			},
			MinRiskRewardRatio:     2.5,
			MinLiquidationDistance: template.Percentage(20),
		},
		Timing: template.TimingConfig{
			DecisionFrequency: template.Duration{
				Value: 5,
				Unit:  "minutes",
			},
			ShortInterval: template.Duration{
				Value: 3,
				Unit:  "minutes",
			},
			LongInterval: template.Duration{
				Value: 4,
				Unit:  "hours",
			},
			RecentDataPointsShort: 50,
			RecentDataPointsLong:  30,
			FocusRecentPoints:     3,
		},
		Output: template.OutputConfig{
			CoinSymbols:           []string{"BTC", "ETH", "SOL"},
			MaxJustificationChars: 500,
		},
	}

	// Create engine
	engine := template.NewJetEngine(template.JetOptions{
		TemplateDir:     "../../etc/prompts/system",
		DevelopmentMode: false,
	})

	// Load template
	tmpl, err := engine.Load("default.jet")
	if err != nil {
		t.Fatalf("Failed to load system prompt template: %v", err)
	}

	// Render
	result, err := engine.Render(tmpl, data)
	if err != nil {
		t.Fatalf("Failed to render system prompt: %v", err)
	}

	// Verify output
	if len(result) == 0 {
		t.Error("Rendered system prompt is empty")
	}

	t.Logf("System prompt length: %d characters", len(result))

	// Check for key content
	expectedParts := []string{
		"GPT-4-Turbo",
		"Hyperliquid",
		"$10.00K",
		"1x to 20x",
		"5 minutes",
		"BTC", "ETH", "SOL",
	}

	for _, part := range expectedParts {
		if !containsString(result, part) {
			t.Errorf("System prompt missing expected part: %q", part)
		}
	}
}

// TestUserPromptRender tests rendering the user prompt template.
func TestUserPromptRender(t *testing.T) {
	// Create sample data
	data := template.UserPromptData{
		Session: template.SessionInfo{
			MinutesElapsed: 120,
		},
		Timeframes: template.TimeframeConfig{
			ShortIntervalMinutes: 3,
			LongIntervalHours:    4,
		},
		Coins: []template.CoinData{
			{
				Symbol: "BTC",
				Current: template.CurrentSnapshot{
					Price: 45000.00,
					EMA20: 44800.00,
					MACD:  150.50,
					RSI7:  65.5,
				},
				Short: template.TimeSeriesData{
					Prices: []float64{44800, 44900, 45000},
					EMA20:  []float64{44700, 44750, 44800},
					MACD:   []float64{145, 148, 150.5},
					RSI7:   []float64{63, 64, 65.5},
					RSI14:  []float64{58, 59, 60},
				},
				Long: template.TimeSeriesData{
					EMA20: []float64{44500, 44600, 44700},
					EMA50: []float64{44000, 44100, 44200},
					ATR3:  []float64{800, 810, 820},
					ATR14: []float64{750, 760, 770},
					MACD:  []float64{140, 145, 150},
					RSI14: []float64{55, 57, 59},
				},
				Futures: template.FuturesMetrics{
					OpenInterest: template.OpenInterestData{
						Latest:  850000000,
						Average: 800000000,
					},
					FundingRate:   0.0001,
					VolumeCurrent: 1500000000,
					VolumeAverage: 1200000000,
				},
			},
		},
		Account: template.AccountInfo{
			Performance: template.PerformanceMetrics{
				ReturnPct:   5.25,
				SharpeRatio: 1.8,
			},
			Status: template.AccountStatus{
				CashAvailable: 8500.00,
				AccountValue:  10500.00,
			},
		},
		Positions: []template.PositionData{
			{
				Symbol:           "BTC",
				Quantity:         0.1,
				EntryPrice:       44000.00,
				CurrentPrice:     45000.00,
				LiquidationPrice: 40000.00,
				UnrealizedPnL:    100.00,
				Leverage:         5,
				ExitPlan: template.ExitPlan{
					ProfitTarget:          48000.00,
					StopLoss:              43000.00,
					InvalidationCondition: "BTC breaks below $42000",
				},
				Confidence:  0.75,
				RiskUSD:     100.00,
				NotionalUSD: 4500.00,
			},
		},
	}

	// Create engine
	engine := template.NewJetEngine(template.JetOptions{
		TemplateDir:     "../../etc/prompts/user",
		DevelopmentMode: false,
	})

	// Load template
	tmpl, err := engine.Load("default.jet")
	if err != nil {
		t.Fatalf("Failed to load user prompt template: %v", err)
	}

	// Render
	result, err := engine.Render(tmpl, data)
	if err != nil {
		t.Fatalf("Failed to render user prompt: %v", err)
	}

	// Verify output
	if len(result) == 0 {
		t.Error("Rendered user prompt is empty")
	}

	t.Logf("User prompt length: %d characters", len(result))

	// Output full result for debugging
	t.Logf("Full output:\n%s", result)

	// Check for key content
	expectedParts := []string{
		"120 minutes",
		"BTC",
		"45000", // without .00 since Jet formats without trailing zeros for whole numbers
		"44800.00, 44900.00, 45000.00",
		"$8.50K",
		"$10.50K",
		"5.25%",
		"1.8",
		"0.1",
		"44000", // entry price
		"48000", // profit target
		"43000", // stop loss
	}

	for _, part := range expectedParts {
		if !containsString(result, part) {
			t.Errorf("User prompt missing expected part: %q", part)
		}
	}
}

// TestCombinedPromptRendering tests rendering both system and user prompts.
func TestCombinedPromptRendering(t *testing.T) {
	// System prompt data
	systemData := template.SystemPromptData{
		Model: template.ModelConfig{Name: "GPT-4"},
		Market: template.MarketConfig{
			Exchange:        "Hyperliquid",
			AssetUniverse:   "BTC, ETH",
			StartingCapital: 5000.00,
			Leverage:        template.Range{Min: 1, Max: 10},
			TradingFee:      template.Range{Min: 0.02, Max: 0.05},
			Slippage:        template.Range{Min: 0.1, Max: 0.3},
		},
		Risk: template.RiskConfig{
			MaxLossPerTrade:    template.Range{Min: 1, Max: 2},
			MinRiskRewardRatio: 2.0,
		},
		Timing: template.TimingConfig{
			DecisionFrequency: template.Duration{Value: 5, Unit: "minutes"},
			ShortInterval:     template.Duration{Value: 3, Unit: "minutes"},
			LongInterval:      template.Duration{Value: 4, Unit: "hours"},
		},
		Output: template.OutputConfig{
			CoinSymbols:           []string{"BTC", "ETH"},
			MaxJustificationChars: 300,
		},
	}

	// User prompt data
	userData := template.UserPromptData{
		Session:    template.SessionInfo{MinutesElapsed: 60},
		Timeframes: template.TimeframeConfig{ShortIntervalMinutes: 3, LongIntervalHours: 4},
		Coins: []template.CoinData{
			{
				Symbol: "BTC",
				Current: template.CurrentSnapshot{
					Price: 45000, EMA20: 44800, MACD: 150, RSI7: 65,
				},
				Short: template.TimeSeriesData{
					Prices: []float64{44800, 45000},
					EMA20:  []float64{44700, 44800},
				},
				Futures: template.FuturesMetrics{
					OpenInterest:  template.OpenInterestData{Latest: 850000000, Average: 800000000},
					FundingRate:   0.0001,
					VolumeCurrent: 1500000000,
					VolumeAverage: 1200000000,
				},
			},
		},
		Account: template.AccountInfo{
			Performance: template.PerformanceMetrics{ReturnPct: 3.5, SharpeRatio: 1.2},
			Status:      template.AccountStatus{CashAvailable: 4800, AccountValue: 5175},
		},
		Positions: []template.PositionData{},
	}

	// Render both prompts
	systemEngine := template.NewJetEngine(template.JetOptions{
		TemplateDir: "../../etc/prompts/system",
	})
	systemTmpl, err := systemEngine.Load("default.jet")
	if err != nil {
		t.Fatalf("Failed to load system template: %v", err)
	}

	userEngine := template.NewJetEngine(template.JetOptions{
		TemplateDir: "../../etc/prompts/user",
	})
	userTmpl, err := userEngine.Load("default.jet")
	if err != nil {
		t.Fatalf("Failed to load user template: %v", err)
	}

	systemResult, err := systemEngine.Render(systemTmpl, systemData)
	if err != nil {
		t.Fatalf("Failed to render system prompt: %v", err)
	}

	userResult, err := userEngine.Render(userTmpl, userData)
	if err != nil {
		t.Fatalf("Failed to render user prompt: %v", err)
	}

	// Combined prompt
	combined := systemResult + "\n\n---\n\n" + userResult

	t.Logf("Combined prompt length: %d characters", len(combined))
	t.Logf("System portion: %d characters", len(systemResult))
	t.Logf("User portion: %d characters", len(userResult))

	// Verify both parts are present
	if !containsString(combined, "GPT-4") {
		t.Error("Combined prompt missing system content")
	}
	if !containsString(combined, "60 minutes") {
		t.Error("Combined prompt missing user content")
	}
}

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
