package repository

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	ent "main/internal/indicator/domain/entities"
	"main/internal/indicator/domain/ports"
)

var (
	rows  pgx.Rows
	query string
)

// postgresqlRepo Struct
type postgresqlRepo struct {
	db *pgxpool.Pool
}

// NewPostgresqlRepository Indicator Domain postgresql repository constructor
func NewPostgresqlRepository(db *pgxpool.Pool) ports.IPostgresqlRepository {
	return &postgresqlRepo{db: db}
}

// WriteResult Save result of indicator calculation
func (r *postgresqlRepo) WriteResult(ctx context.Context, req_body *ent.IndicatorCalcResponse) (err error) {
	json_data, _ := json.Marshal(req_body)

	query := `INSERT INTO response_collector (pair, response_data) VALUES ($1, $2)`
	_, err = r.db.Query(ctx, query, "BTCUSDTT", json_data)
	if err != nil {
		return errors.Wrap(err, "indicatorPostgresqlRepo.WriteResult")
	}

	return nil
}

// ReadIndicatorIndexByInterval Get Record of Booster Indexes
func (r *postgresqlRepo) ReadIndicatorIndexByInterval(ctx context.Context, intervals string) (result []ent.IndicatorIndex, err error) {
	query = `SELECT id, exchange, indicator_id, func_name, intervals, asset_id, asset_symbol, channel_tag FROM indicators_index WHERE intervals=$1 AND total_indicator>0 AND is_locked=false`
	rows, err = r.db.Query(ctx, query, intervals)
	if err != nil {
		println(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		dat := ent.IndicatorIndex{}
		err2 := rows.Scan(&dat.Id, &dat.Exchange, &dat.IndicatorId, &dat.FuncName, &dat.Intervals, &dat.AssetsId, &dat.AssetsSymbol, &dat.ChannelTag)
		if err2 != nil {
			println(err.Error())
		}

		result = append(result, dat)
	}

	return
}

// ChangeIndicatorIndexStatus change latest read status of indicator row
func (r *postgresqlRepo) ChangeIndicatorIndexStatus(ctx context.Context, id int64) (err error) {

	query = `UPDATE indicators_index SET calculated_at=NOW() WHERE id=$1`
	_, err = r.db.Exec(ctx, query, id)

	return
}

// ChangeLockStatusByRow change latest read status of indicator row
func (r *postgresqlRepo) ChangeLockStatusByRow(ctx context.Context, id int64) (err error) {

	query = `UPDATE indicators_index SET is_locked=true WHERE id=$1`
	_, err = r.db.Exec(ctx, query, id)

	return
}

// ReleaseAllLocks change latest read status of indicator row
func (r *postgresqlRepo) ReleaseAllLocks(ctx context.Context) (err error) {

	query = `UPDATE indicators_index SET is_locked=false`
	_, err = r.db.Exec(ctx, query)

	return
}

// UpdateIndicatorResult change latest read status of indicator row
func (r *postgresqlRepo) UpdateIndicatorResult(ctx context.Context, row ent.IndicatorIndex, res ent.IndicatorCalcResponse) (err error) {

	query = `UPDATE booster_sub SET signal=$1, result=$2 WHERE exchange=$3 AND indicator_id=$4 AND asset_id=$5 AND intervals=$6`
	_, err = r.db.Exec(ctx, query, res.Signal, res.Result, row.Exchange, row.IndicatorId, row.AssetsId, row.Intervals)

	return
}

// SaveBoosterStatistics Save result of booster statistics calculation
func (r *postgresqlRepo) SaveBoosterStatistics(ctx context.Context, intervals string, asset_id int64, asset_symbol string, req_body ent.BoosterStatisticsResponse) (err error) {
	//json_data, _ := json.Marshal(req_body)

	query := `INSERT INTO booster_index AS bi (intervals, asset_id, asset_symbol, total_pair, line, fibo, mmath) VALUES ($1, $2, $3, 1, $4, $5, $6) ON CONFLICT (intervals, asset_id) DO UPDATE SET total_pair=(bi.total_pair+1), line=$4, fibo=$5, mmath=$6`
	_, err = r.db.Query(ctx, query, intervals, asset_id, asset_symbol, req_body.Line, req_body.Fibonacci, req_body.MMath)
	if err != nil {
		return errors.Wrap(err, "indicatorPostgresqlRepo.SaveBoosterStatistics")
	}

	return err
}
