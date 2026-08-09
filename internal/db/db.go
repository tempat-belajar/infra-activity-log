package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/akr/infra-activity-log/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

func (d *DB) Close() {
	d.Pool.Close()
}

type ListFilter struct {
	DateFrom string
	DateTo   string
	PIC      string
	Status   string
	Category string
	Search   string
}

func (d *DB) ListLogs(ctx context.Context, f ListFilter) ([]models.ActivityLog, error) {
	q := `SELECT id, tanggal, job_title, pic, application, label,
		old_value_text, old_value_image_url, new_value_text, new_value_image_url,
		status, category, created_at, updated_at
		FROM activity_logs WHERE 1=1`
	args := []interface{}{}
	argN := 1

	if f.DateFrom != "" {
		q += fmt.Sprintf(" AND tanggal >= $%d", argN)
		args = append(args, f.DateFrom)
		argN++
	}
	if f.DateTo != "" {
		q += fmt.Sprintf(" AND tanggal <= $%d", argN)
		args = append(args, f.DateTo)
		argN++
	}
	if f.PIC != "" {
		q += fmt.Sprintf(" AND pic = $%d", argN)
		args = append(args, f.PIC)
		argN++
	}
	if f.Status != "" {
		q += fmt.Sprintf(" AND status = $%d", argN)
		args = append(args, f.Status)
		argN++
	}
	if f.Category != "" {
		q += fmt.Sprintf(" AND category = $%d", argN)
		args = append(args, f.Category)
		argN++
	}
	if f.Search != "" {
		q += fmt.Sprintf(" AND (label ILIKE $%d OR application ILIKE $%d OR job_title ILIKE $%d)", argN, argN, argN)
		args = append(args, "%"+f.Search+"%")
		argN++
	}
	q += " ORDER BY tanggal DESC, id DESC"

	rows, err := d.Pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []models.ActivityLog
	for rows.Next() {
		var l models.ActivityLog
		if err := rows.Scan(&l.ID, &l.Tanggal, &l.JobTitle, &l.PIC, &l.Application, &l.Label,
			&l.OldValueText, &l.OldValueImageURL, &l.NewValueText, &l.NewValueImageURL,
			&l.Status, &l.Category, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []models.ActivityLog{}
	}
	return logs, rows.Err()
}

func (d *DB) GetLog(ctx context.Context, id int) (*models.ActivityLog, error) {
	q := `SELECT id, tanggal, job_title, pic, application, label,
		old_value_text, old_value_image_url, new_value_text, new_value_image_url,
		status, category, created_at, updated_at
		FROM activity_logs WHERE id = $1`
	var l models.ActivityLog
	err := d.Pool.QueryRow(ctx, q, id).Scan(&l.ID, &l.Tanggal, &l.JobTitle, &l.PIC, &l.Application, &l.Label,
		&l.OldValueText, &l.OldValueImageURL, &l.NewValueText, &l.NewValueImageURL,
		&l.Status, &l.Category, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (d *DB) CreateLog(ctx context.Context, l models.ActivityLog) (int, error) {
	q := `INSERT INTO activity_logs
		(tanggal, job_title, pic, application, label, old_value_text, old_value_image_url,
		new_value_text, new_value_image_url, status, category)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id`
	var id int
	err := d.Pool.QueryRow(ctx, q, l.Tanggal, l.JobTitle, l.PIC, l.Application, l.Label,
		l.OldValueText, l.OldValueImageURL, l.NewValueText, l.NewValueImageURL,
		l.Status, l.Category).Scan(&id)
	return id, err
}

func (d *DB) UpdateLog(ctx context.Context, id int, l models.ActivityLog, keepOldImage, keepNewImage bool) error {
	setParts := []string{
		"tanggal=$2", "job_title=$3", "pic=$4", "application=$5", "label=$6",
		"old_value_text=$7", "new_value_text=$8", "status=$9", "category=$10", "updated_at=now()",
	}
	args := []interface{}{id, l.Tanggal, l.JobTitle, l.PIC, l.Application, l.Label,
		l.OldValueText, l.NewValueText, l.Status, l.Category}
	argN := 11

	if !keepOldImage {
		setParts = append(setParts, fmt.Sprintf("old_value_image_url=$%d", argN))
		args = append(args, l.OldValueImageURL)
		argN++
	}
	if !keepNewImage {
		setParts = append(setParts, fmt.Sprintf("new_value_image_url=$%d", argN))
		args = append(args, l.NewValueImageURL)
		argN++
	}

	q := fmt.Sprintf("UPDATE activity_logs SET %s WHERE id=$1", strings.Join(setParts, ", "))
	_, err := d.Pool.Exec(ctx, q, args...)
	return err
}

func (d *DB) DeleteLog(ctx context.Context, id int) error {
	_, err := d.Pool.Exec(ctx, "DELETE FROM activity_logs WHERE id=$1", id)
	return err
}
