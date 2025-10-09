package model

import "time"

type DailySalesReport struct {
	Day	time.Time
	Revenue	float64
	Orders	int
	UniqueCustomers	int
	AOV	float64 
}

type TopProductReport struct{
	ProductId int
	ProductName string
	Quantity int
	GMV float64
}

type orderBystatusReport struct{
	Status string
	Count int
	Revenue float64
}