package dbx

import (
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

func QueryBatchSelect(b *pgx.Batch, sb squirrel.SelectBuilder) error {
	query, args, err := sb.ToSql()
	if err != nil {
		return err
	}
	b.Queue(query, args...)
	return nil
}
