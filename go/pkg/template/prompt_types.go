package template

// SystemPromptData contains all data for rendering the system prompt template.
type SystemPromptData struct {
	Model  ModelConfig  `json:"model" doc:"Model configuration"`
	Market MarketConfig `json:"market" doc:"Market and trading environment configuration"`
	Risk   RiskConfig   `json:"risk" doc:"Risk management parameters"`
	Timing TimingConfig `json:"timing" doc:"Timing and frequency settings"`
	Output OutputConfig `json:"output" doc:"Output format configuration"`
}

// ModelConfig contains model identification and setup.
type ModelConfig struct {
	Name string `json:"name" doc:"Model name/designation" example:"GPT-4"`
}

// MarketConfig contains market and trading environment parameters.
type MarketConfig struct {
	Exchange                 string     `json:"exchange" doc:"Exchange name" example:"Hyperliquid"`
	AssetUniverse            string     `json:"asset_universe" doc:"Description of tradeable assets" example:"BTC, ETH, SOL"`
	StartingCapital          float64    `json:"starting_capital" doc:"Initial capital in USD" example:"10000"`
	MarketHours              string     `json:"market_hours" doc:"Trading hours" example:"24/7"`
	ContractType             string     `json:"contract_type" doc:"Type of contracts" example:"Perpetual futures"`
	Leverage                 Range      `json:"leverage" doc:"Allowed leverage range" example:"{\"min\":1,\"max\":20}"`
	TradingFee               Range      `json:"trading_fee" doc:"Trading fee percentage range" example:"{\"min\":0.02,\"max\":0.05}"`
	Slippage                 Range      `json:"slippage" doc:"Expected slippage percentage range" example:"{\"min\":0.1,\"max\":0.5}"`
	MinPositionSize          float64    `json:"min_position_size" doc:"Minimum position size in USD" example:"100"`
	MaxPositionConcentration Percentage `json:"max_position_concentration" doc:"Maximum % of capital in single position" example:"30"`
}

// RiskConfig contains risk management parameters.
type RiskConfig struct {
	MaxLossPerTrade        Range      `json:"max_loss_per_trade" doc:"Acceptable loss per trade (% of account)" example:"{\"min\":1,\"max\":3}"`
	MinRiskRewardRatio     float64    `json:"min_risk_reward_ratio" doc:"Minimum reward-to-risk ratio" example:"2.5"`
	MinLiquidationDistance Percentage `json:"min_liquidation_distance" doc:"Minimum distance from liquidation (%)" example:"20"`
}

// TimingConfig contains timing and frequency settings.
type TimingConfig struct {
	DecisionFrequency     Duration `json:"decision_frequency" doc:"How often to make decisions" example:"{\"value\":5,\"unit\":\"minutes\"}"`
	ShortInterval         Duration `json:"short_interval" doc:"Short-term data interval" example:"{\"value\":3,\"unit\":\"minutes\"}"`
	LongInterval          Duration `json:"long_interval" doc:"Long-term data interval" example:"{\"value\":4,\"unit\":\"hours\"}"`
	RecentDataPointsShort int      `json:"recent_data_points_short" doc:"Number of recent short-interval data points" example:"50"`
	RecentDataPointsLong  int      `json:"recent_data_points_long" doc:"Number of recent long-interval data points" example:"30"`
	FocusRecentPoints     int      `json:"focus_recent_points" doc:"Number of most recent points to focus on" example:"3"`
}

// OutputConfig contains output format configuration.
type OutputConfig struct {
	CoinSymbols           []string `json:"coin_symbols" doc:"List of tradeable coin symbols" example:"[\"BTC\",\"ETH\",\"SOL\"]"`
	MaxJustificationChars int      `json:"max_justification_chars" doc:"Maximum characters in trade justification" example:"500"`
}

// UserPromptData contains all data for rendering the user prompt template.
type UserPromptData struct {
	Session    SessionInfo     `json:"session" doc:"Trading session information"`
	Timeframes TimeframeConfig `json:"timeframes" doc:"Timeframe configuration"`
	Coins      []CoinData      `json:"coins" doc:"Market data for all coins"`
	Account    AccountInfo     `json:"account" doc:"Account status and performance"`
	Positions  []PositionData  `json:"positions" doc:"Current open positions"`
}

// SessionInfo contains trading session information.
type SessionInfo struct {
	MinutesElapsed int `json:"minutes_elapsed" doc:"Minutes since trading started" example:"120"`
}

// TimeframeConfig contains timeframe settings.
type TimeframeConfig struct {
	ShortIntervalMinutes int `json:"short_interval_minutes" doc:"Short-term interval in minutes" example:"3"`
	LongIntervalHours    int `json:"long_interval_hours" doc:"Long-term interval in hours" example:"4"`
}

// CoinData contains comprehensive market data for a single coin.
type CoinData struct {
	Symbol  string          `json:"symbol" doc:"Coin symbol" example:"BTC"`
	Current CurrentSnapshot `json:"current" doc:"Current market snapshot"`
	Short   TimeSeriesData  `json:"short" doc:"Short-term time series data"`
	Long    TimeSeriesData  `json:"long" doc:"Long-term time series data"`
	Futures FuturesMetrics  `json:"futures" doc:"Perpetual futures specific metrics"`
}

// CurrentSnapshot contains current market state.
type CurrentSnapshot struct {
	Price float64 `json:"price" doc:"Current price" example:"45000.00"`
	EMA20 float64 `json:"ema20" doc:"20-period EMA" example:"44800.00"`
	MACD  float64 `json:"macd" doc:"MACD indicator" example:"150.50"`
	RSI7  float64 `json:"rsi7" doc:"7-period RSI" example:"65.5"`
}

// TimeSeriesData contains time series indicators.
type TimeSeriesData struct {
	Prices []float64 `json:"prices" doc:"Price series (oldest to newest)" example:"[45000, 45100, 45200]"`
	EMA20  []float64 `json:"ema20" doc:"20-period EMA series" example:"[44800, 44850, 44900]"`
	EMA50  []float64 `json:"ema50" doc:"50-period EMA series (long-term only)" example:"[44500, 44550, 44600]"`
	MACD   []float64 `json:"macd" doc:"MACD series" example:"[150, 155, 160]"`
	RSI7   []float64 `json:"rsi7" doc:"7-period RSI series" example:"[63, 64, 65]"`
	RSI14  []float64 `json:"rsi14" doc:"14-period RSI series" example:"[58, 59, 60]"`
	ATR3   []float64 `json:"atr3" doc:"3-period ATR series (long-term only)" example:"[800, 810, 820]"`
	ATR14  []float64 `json:"atr14" doc:"14-period ATR series (long-term only)" example:"[750, 760, 770]"`
}

// FuturesMetrics contains perpetual futures specific data.
type FuturesMetrics struct {
	OpenInterest  OpenInterestData `json:"open_interest" doc:"Open interest data"`
	FundingRate   float64          `json:"funding_rate" doc:"Current funding rate" example:"0.0001"`
	VolumeCurrent float64          `json:"volume_current" doc:"Current volume" example:"1500000000"`
	VolumeAverage float64          `json:"volume_average" doc:"Average volume" example:"1200000000"`
}

// OpenInterestData contains open interest information.
type OpenInterestData struct {
	Latest  float64 `json:"latest" doc:"Latest open interest" example:"850000000"`
	Average float64 `json:"average" doc:"Average open interest" example:"800000000"`
}

// AccountInfo contains account status and performance.
type AccountInfo struct {
	Performance PerformanceMetrics `json:"performance" doc:"Performance metrics"`
	Status      AccountStatus      `json:"status" doc:"Current account status"`
}

// PerformanceMetrics contains performance data.
type PerformanceMetrics struct {
	ReturnPct   float64 `json:"return_pct" doc:"Total return percentage" example:"5.25"`
	SharpeRatio float64 `json:"sharpe_ratio" doc:"Sharpe ratio" example:"1.8"`
}

// AccountStatus contains current account state.
type AccountStatus struct {
	CashAvailable float64 `json:"cash_available" doc:"Available cash in USD" example:"8500.00"`
	AccountValue  float64 `json:"account_value" doc:"Total account value in USD" example:"10500.00"`
}

// PositionData contains information about an open position.
type PositionData struct {
	Symbol           string   `json:"symbol" doc:"Position symbol" example:"BTC"`
	Quantity         float64  `json:"quantity" doc:"Position quantity" example:"0.1"`
	EntryPrice       float64  `json:"entry_price" doc:"Entry price" example:"45000.00"`
	CurrentPrice     float64  `json:"current_price" doc:"Current price" example:"46000.00"`
	LiquidationPrice float64  `json:"liquidation_price" doc:"Liquidation price" example:"40000.00"`
	UnrealizedPnL    float64  `json:"unrealized_pnl" doc:"Unrealized profit/loss" example:"100.00"`
	Leverage         int      `json:"leverage" doc:"Position leverage" example:"5"`
	ExitPlan         ExitPlan `json:"exit_plan" doc:"Exit strategy"`
	Confidence       float64  `json:"confidence" doc:"Trade confidence (0-1)" example:"0.75"`
	RiskUSD          float64  `json:"risk_usd" doc:"Risk amount in USD" example:"150.00"`
	NotionalUSD      float64  `json:"notional_usd" doc:"Notional value in USD" example:"4500.00"`
}

// ExitPlan contains exit strategy parameters.
type ExitPlan struct {
	ProfitTarget          float64 `json:"profit_target" doc:"Take profit price" example:"48000.00"`
	StopLoss              float64 `json:"stop_loss" doc:"Stop loss price" example:"44000.00"`
	InvalidationCondition string  `json:"invalidation_condition" doc:"Condition that invalidates the trade" example:"BTC breaks below $43000"`
}
