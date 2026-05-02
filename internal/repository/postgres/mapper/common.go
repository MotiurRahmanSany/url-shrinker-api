package mapper

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgUUID(uuid string) (pgtype.UUID, error) {
	var pgUUID pgtype.UUID
	err := pgUUID.Scan(uuid)
	return pgUUID, err
}

func FromPgUUID(pgUUID pgtype.UUID) string {
	return pgUUID.String()
}

func ToTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{
		Time:  *t,
		Valid: true,
	}
}

func FromTimestamptz(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

func ToPgInt4(i *int32) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{
		Int32: *i,
		Valid: true,
	}
}

func FromPgInt4(i pgtype.Int4) *int32 {
	if !i.Valid {
		return nil
	}
	return &i.Int32
}

func ToPgText(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{
		String: *s,
		Valid:  true,
	}
}

func FromPgText(t pgtype.Text) *string {
	if !t.Valid {
		return nil
	}
	return &t.String
}
