package accrualprocessor

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/shopspring/decimal"
	"github.com/thorgnir-go-study/yp-diploma/internal/entity"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/accrual"
	"github.com/thorgnir-go-study/yp-diploma/internal/usecase/order"
	"time"

	"github.com/rs/zerolog/log"
)

type Service struct {
	repo      Repository
	accrualUC accrual.UseCase
	orderUC   order.UseCase
	jobsCh    chan<- entity.ProcessingTask
	scheduler *gocron.Scheduler
}

func NewService(repo Repository, accrualUC accrual.UseCase, orderUC order.UseCase) *Service {
	s := &Service{repo: repo, accrualUC: accrualUC, orderUC: orderUC}
	// TODO: притащить контекст (graceful shutdown), кол-во воркеров, интервалы крона вынести в конфиг
	s.jobsCh = s.startWorkers(context.Background(), 5)
	s.scheduler = gocron.NewScheduler(time.Local)
	return s
}

func (s *Service) Start() error {
	if s.scheduler.IsRunning() {
		return nil
	}
	s.scheduler.Clear()
	_, _ = s.scheduler.Every(3).Seconds().Do(func() { s.scheduleNewOrders(context.Background()) })
	_, _ = s.scheduler.Every(3).Seconds().Do(func() { s.getJobsToRun(context.Background()) })
	_, _ = s.scheduler.Every(10).Seconds().Do(func() { s.cleanProcessedTasks(context.Background()) })
	s.scheduler.StartAsync()
	return nil
}

func (s *Service) scheduleNewOrders(ctx context.Context) {
	newOrders, err := s.orderUC.GetNewOrders(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error while getting new orders")
		return
	}
	if len(newOrders) > 0 {
		var newTasks = make([]*entity.ProcessingTask, len(newOrders))
		for i := range newOrders {
			o := *newOrders[i]
			task, err := entity.NewProcessingTask(o.ID, o.Number, time.Now())
			if err != nil {
				log.Error().Err(err).Msg("Error while creating new processing task")
				continue
			}
			newTasks[i] = task
		}
		err = s.repo.CreateTasks(ctx, newTasks)
		if err != nil {
			log.Error().Err(err).Msg("Error while creating new processing tasks in repo")
			return
		}
	}

}

func (s *Service) getJobsToRun(ctx context.Context) {
	tasks, err := s.repo.GetTasksToRun(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error while getting tasks to run")
		return
	}
	if len(tasks) > 0 {
		for _, task := range tasks {
			s.jobsCh <- *task
		}
	}
}

func (s *Service) cleanProcessedTasks(ctx context.Context) {
	err := s.repo.CleanProcessedTasks(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Error while cleaning processed tasks")
	}
}

func (s *Service) startWorkers(ctx context.Context, workersCount int) chan<- entity.ProcessingTask {
	ch := make(chan entity.ProcessingTask, workersCount*2)
	for i := 0; i < workersCount; i++ {
		workerID := fmt.Sprintf("AccrualProcessingWorkerId#%d", i+1)
		go func() {
			log.Info().Str("worker", workerID).Msg("starting accruals processing worker")
			for {
				select {
				case <-ctx.Done():
					log.Info().Str("worker", workerID).Msg("stopping accruals processing worker")
				case req := <-ch:
					go func() {
						innerCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
						defer cancel()
						err := s.repo.SetTaskStatus(ctx, req.ID, entity.ProcessingTaskStatusProcessing)
						if err != nil {
							log.Error().Err(err).
								Str("worker", workerID).
								Str("orderID", req.OrderID.String()).
								Str("status", entity.ProcessingTaskStatusProcessing.String()).
								Msg("error updating task status")
							return
						}

						result, err := s.accrualUC.GetAccrual(innerCtx, req.OrderNumber)
						if err != nil {
							log.Error().
								Err(err).
								Str("worker", workerID).
								Str("orderNumber", req.OrderNumber.String()).
								Msg("error while getting order accrual")

							_ = s.repo.RescheduleTask(ctx, req.ID, calcNextRun())
							return
						}
						log.Info().
							Str("workerId", workerID).
							Str("orderNumber", req.OrderNumber.String()).
							Str("status", result.Status.String()).
							Str("accrual", result.Accrual.String()).
							Msg("Got accrual response from external system")
						switch result.Status {
						case entity.AccrualOrderStatusRegistered, entity.AccrualOrderStatusProcessing:
							{
								_ = s.repo.RescheduleTask(ctx, req.ID, calcNextRun())
								break
							}
						case entity.AccrualOrderStatusInvalid:
							{
								_ = s.orderUC.SetOrderStatus(ctx, req.OrderID, entity.OrderStatusInvalid)
							}
						case entity.AccrualOrderStatusProcessed:
							{
								acc := decimal.NewNullDecimal(decimal.NewFromInt(0))
								if result.Accrual != nil {
									acc = decimal.NewNullDecimal(*result.Accrual)
								}
								_ = s.orderUC.SetOrderAccrualAndStatus(ctx, req.OrderID, acc, entity.OrderStatusProcessed)
								_ = s.repo.SetTaskStatus(ctx, req.ID, entity.ProcessingTaskStatusProcessed)
							}

						}
					}()

				}
			}
		}()
	}
	return ch
}

func calcNextRun() time.Time {
	return time.Now().Add(5 * time.Second)
}
