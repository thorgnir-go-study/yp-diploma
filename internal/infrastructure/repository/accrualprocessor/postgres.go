package accrualprocessor

import (
	context "context"
	"errors"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"time"
)

type PostgresAccrualProcessorRepository struct {
	dbpool *pgxpool.Pool
}

type dbEntity struct {
	ID          uuid.UUID                   `db:"id"`
	OrderID     uuid.UUID                   `db:"order_id"`
	ToRunAt     time.Time                   `db:"to_run_at"`
	Status      entity.ProcessingTaskStatus `db:"status"`
	UpdatedAt   time.Time                   `db:"updated_at"`
	OrderNumber *string                     `db:"order_number"`
}

func NewPostgresAccrualProcessorRepository(dbpool *pgxpool.Pool) *PostgresAccrualProcessorRepository {
	return &PostgresAccrualProcessorRepository{dbpool: dbpool}
}

func (p PostgresAccrualProcessorRepository) GetTasksToRun(ctx context.Context) ([]*entity.ProcessingTask, error) {
	log.Debug().Msg("Getting tasks to run")
	var tasks []*dbEntity
	if err := pgxscan.Select(ctx, p.dbpool, &tasks, `
select p.id, p.order_id, p.to_run_at, p.status, p.updated_at, o.order_number as order_number
from gophermart."processing_task" p
left join gophermart."order" o on p.order_id = o.id
where to_run_at <= $1
`, time.Now()); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		log.Error().Err(err).Msg("Error while getting tasks to run")
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}

	result := make([]*entity.ProcessingTask, len(tasks))
	for i := range tasks {
		t := tasks[i]
		orderNumber, _ := entity.StringToOrderNumber(*t.OrderNumber)
		result[i] = &entity.ProcessingTask{
			ID:          t.ID,
			OrderID:     t.OrderID,
			OrderNumber: orderNumber,
			ToRunAt:     t.ToRunAt,
			Status:      t.Status,
			UpdatedAt:   t.UpdatedAt,
		}
	}

	return result, nil
}

func (p PostgresAccrualProcessorRepository) CreateTasks(ctx context.Context, tasks []*entity.ProcessingTask) error {
	tx, err := p.dbpool.Begin(ctx)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback(ctx) //nolint:errcheck

	for _, task := range tasks {
		if _, err = tx.Exec(ctx, `
insert into gophermart."processing_task" (id, order_id, to_run_at, status, updated_at)
values ($1, $2, $3, $4, $5)
`, task.ID, task.OrderID, task.ToRunAt, task.Status, task.UpdatedAt); err != nil {
			log.Error().Err(err).Msg("Error while inserting processing task")
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error while commiting processing tasks creating transaction")
		return err
	}
	return nil
}

func (p PostgresAccrualProcessorRepository) SetTaskStatus(ctx context.Context, taskID entity.ID, status entity.ProcessingTaskStatus) error {
	if _, err := p.dbpool.Exec(ctx, `
update gophermart."processing_task"
set status = $2, updated_at = $3
where id = $1;
`, taskID, status, time.Now()); err != nil {
		return err
	}
	log.Debug().Str("taskID", taskID.String()).Str("status", status.String()).Msg("Set task status")
	return nil
}

func (p PostgresAccrualProcessorRepository) RescheduleTask(ctx context.Context, taskID entity.ID, nextRun time.Time) error {
	log.Debug().Str("taskID", taskID.String()).Str("nextRun", nextRun.String()).Msg("Reschedule")
	if _, err := p.dbpool.Exec(ctx, `
update gophermart."processing_task"
set status = $2, to_run_at = $3, updated_at = $4
where id = $1;
`, taskID, entity.ProcessingTaskStatusScheduled, nextRun, time.Now()); err != nil {
		return err
	}
	return nil
}

func (p PostgresAccrualProcessorRepository) CleanProcessedTasks(ctx context.Context) error {
	if _, err := p.dbpool.Exec(ctx, `
DELETE FROM gophermart."processing_task"
where status = $1;
`, entity.ProcessingTaskStatusProcessed); err != nil {
		return err
	}
	return nil
}
