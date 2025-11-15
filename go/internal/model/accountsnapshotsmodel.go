package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AccountSnapshotsModel = (*customAccountSnapshotsModel)(nil)

type (
	// AccountSnapshotsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountSnapshotsModel.
	AccountSnapshotsModel interface {
		accountSnapshotsModel
	}

	customAccountSnapshotsModel struct {
		*defaultAccountSnapshotsModel
	}
)

// NewAccountSnapshotsModel returns a model for the database table.
func NewAccountSnapshotsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AccountSnapshotsModel {
	return &customAccountSnapshotsModel{
		defaultAccountSnapshotsModel: newAccountSnapshotsModel(conn, c, opts...),
	}
}
