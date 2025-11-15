package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ KlinesModel = (*customKlinesModel)(nil)

type (
	// KlinesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customKlinesModel.
	KlinesModel interface {
		klinesModel
	}

	customKlinesModel struct {
		*defaultKlinesModel
	}
)

// NewKlinesModel returns a model for the database table.
func NewKlinesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) KlinesModel {
	return &customKlinesModel{
		defaultKlinesModel: newKlinesModel(conn, c, opts...),
	}
}
