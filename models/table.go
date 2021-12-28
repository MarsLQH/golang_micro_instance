package models

import (
	"database/sql"
	"gorm.io/datatypes"
	"time"
)

type Table struct {
	Id uint32 `json:"id"`

	ColInt32     int32          `json:"col_int_32"`
	ColJson      datatypes.JSON `json:"col_json,omitempty"`
	OpAt         sql.NullTime   `json:"recommend_at,omitempty"`
	RecommendPos datatypes.JSON `json:"recommend_pos,omitempty"`
	ColString    string         `gorm:"column:col_strings" json:"cid,omitempty"`
	CreateAt     time.Time      `json:"create_at"`
}
