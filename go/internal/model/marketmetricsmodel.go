package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MarketMetricsModel = (*customMarketMetricsModel)(nil)

type (
	// MarketMetricsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMarketMetricsModel.
	MarketMetricsModel interface {
		marketMetricsModel
	}

	customMarketMetricsModel struct {
		*defaultMarketMetricsModel
	}
)

// NewMarketMetricsModel returns a model for the database table.
func NewMarketMetricsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MarketMetricsModel {
	return &customMarketMetricsModel{
		defaultMarketMetricsModel: newMarketMetricsModel(conn, c, opts...),
	}
}
