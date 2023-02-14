package entities

import (
	"time"
)

type IType string
type ISignal int
type Trend int

const (
	Falling Trend = -1
	Rising  Trend = 1
)

type IndicatorCalcRequest struct {
	Provider  string `json:"provider" validate:"required"`
	Indicator string `json:"indicator" validate:"required"`
	Pair      string `json:"pair" validate:"required"`
	Interval  string `json:"interval" validate:"required"`
	Limit     int    `json:"limit"`
	Period    int    `json:"period,omitempty"`
	IsTest    bool   `json:"is_test,omitempty"`
}

type Indicator struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Script      string    `json:"script"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	IsOptimized bool      `json:"is_optimized,omitempty"`
}

type IndicatorCalcResponse struct {
	Signal       int         `json:"signal"`
	Signals      interface{} `json:"signals,omitempty"`
	Result       interface{} `json:"results,omitempty"`
	CalculatedAt int64       `json:"calculated_at,omitempty"`
}

type IndicatorIndex struct {
	Id           int64  `json:"id,omitempty"`
	Exchange     string `json:"exchange,omitempty"`
	IndicatorId  int64  `json:"indicator_id,omitempty"`
	FuncName     string `json:"func_name,omitempty"`
	Intervals    string `json:"intervals,omitempty"`
	AssetsId     int64  `json:"assets_id,omitempty"`
	AssetsSymbol string `json:"assets_symbol,omitempty"`
	ChannelTag   string `json:"channel_tag,omitempty"`
}

type BoosterIndex struct {
	Id           int64  `json:"id,omitempty"`
	Exchange     string `json:"exchange,omitempty"`
	IndicatorId  int64  `json:"indicator_id,omitempty"`
	FuncName     string `json:"func_name,omitempty"`
	Intervals    string `json:"intervals,omitempty"`
	AssetsId     int64  `json:"assets_id,omitempty"`
	AssetsSymbol string `json:"assets_symbol,omitempty"`
	ChannelTag   string `json:"channel_tag,omitempty"`
}

type Line struct {
	Rsi         int `json:"rsi"`
	Macd        int `json:"macd"`
	Tds         int `json:"tds"`
	GoldenCross int `json:"golden_cross"`
	SilverCross int `json:"silver_cross"`
	DeathCross  int `json:"death_cross"`
}

type BoosterStatisticsResponse struct {
	Line      Line        `json:"line,omitempty"`
	Fibonacci interface{} `json:"fibonacci,omitempty"`
	MMath     interface{} `json:"m_math,omitempty"`
}
