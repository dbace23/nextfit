package repository

import (
	"context"
	"database/sql"
	"nextfit/internal/model"
)

type ReportRepository struct {
	DB *sql.DB
}

func (reportRepository *ReportRepository) GetSalesSummaryDaily(ctx context.Context) ([]model.DailySalesReport, error) {
	const query = `
	WITH RECURSIVE dates AS (
		SELECT CURDATE() - INTERVAL 29 DAY AS d
		UNION ALL
		SELECT DATE_ADD(d, INTERVAL 1 DAY) FROM dates WHERE d < CURDATE()
	),
	agg AS (
		SELECT
			DATE(o.created_at) AS d,
			SUM(o.paid_amount)        AS revenue,
			COUNT(*)                  AS orders,
			COUNT(DISTINCT o.user_id) AS customers
		FROM orders o
		WHERE o.deleted_at IS NULL
		AND o.created_at >= CURDATE() - INTERVAL 29 DAY
		-- AND o.order_status = 'COMPLETED'
		GROUP BY DATE(o.created_at)
	)
	SELECT
		dates.d,
		COALESCE(agg.revenue, 0),
		COALESCE(agg.orders, 0),
		COALESCE(agg.customers, 0)
	FROM dates
	LEFT JOIN agg ON agg.d = dates.d
	ORDER BY dates.d;
	`
	result, err := reportRepository.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var out []model.DailySalesReport
	for result.Next() {
		var d model.DailySalesReport
		if err := result.Scan(&d.Day, &d.Revenue, &d.Orders, &d.UniqueCustomers); err != nil {
			return nil, err
		}
		if d.Orders > 0 {
			d.AOV = d.Revenue / float64(d.Orders)
		}
		out = append(out, d)
	}
	return out, result.Err()
}

 
func (reportRepository *ReportRepository) GetTopProducts(ctx context.Context) ([]model.TopProductReport, error) {
	const query = `
	SELECT
		p.product_id,
		p.product_name,
		SUM(od.quantity)                   AS qty,
		SUM(od.quantity * p.selling_price) AS gmv
	FROM orders o
	JOIN order_details od ON o.order_id = od.order_id
	JOIN products p       ON p.product_id = od.product_id
	WHERE o.deleted_at IS NULL
	AND od.deleted_at IS NULL
	AND p.deleted_at IS NULL
	AND o.created_at >= CURDATE() - INTERVAL 29 DAY
	AND o.order_status = 'COMPLETED'
	GROUP BY p.product_id, p.product_name
	ORDER BY gmv DESC
	LIMIT 5;
	`
	result, err := reportRepository.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var out []model.TopProductReport
	for result.Next() {
		var r0 model.TopProductReport
		if err :=result.Scan(&r0.ProductId, &r0.ProductName, &r0.Quantity, &r0.GMV); err != nil {
			return nil, err
		}
		out = append(out, r0)
	}
	return out, result.Err()
}

 
func (reportRepository *ReportRepository) GetOrdersByStatus(ctx context.Context) ([]model.OrdersByStatusReport, error) {
	const query = `
	SELECT
		o.order_status,
		COUNT(*)         AS cnt,
		SUM(o.paid_amount) AS revenue
	FROM orders o
	WHERE o.deleted_at IS NULL
	AND o.created_at >= CURDATE() - INTERVAL 29 DAY
	GROUP BY o.order_status
	ORDER BY o.order_status;
	`
	result, err := reportRepository.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer result.Close()
