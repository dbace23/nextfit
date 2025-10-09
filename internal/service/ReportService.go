package service

import (
	"context"
	"nextfit/internal/model"
	"nextfit/internal/repository"
)

type ReportService struct {
	Repo *repository.ReportRepository
}

func NewReportService(repo *repository.ReportRepository) *ReportService {
	return &ReportService{Repo: repo}
}

func (s *ReportService) GetSalesSummaryDaily() ([]model.DailySalesReport, error) {
	return s.Repo.GetSalesSummaryDaily(context.Background())
}

func (s *ReportService) GetTopProducts() ([]model.TopProductReport, error) {
	return s.Repo.GetTopProducts(context.Background())
}

func (s *ReportService) GetOrdersByStatus() ([]model.OrdersByStatusReport, error) {
    return s.Repo.GetOrdersByStatus(context.Background())
}

