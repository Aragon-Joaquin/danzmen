package db

import (
	ty "danzmen/types"
	"database/sql"
)

type DBMonthly_Task struct {
	Id   int
	Name string
}

type DBMonthly_Record struct {
	Year_MonthId   int
	MonthlyId      int
	Completed_At   sql.NullString
	Times_Done     sql.NullFloat64
	Times_Required float64
}

type DBJoin_Monthly struct {
	*DBMonthly_Task
	*DBMonthly_Record
	*DBYear_Month

	ty.MonthlyTasksCfg
}

type DBLong_Tasks struct {
	Id             int
	Name           string
	Expires_in     sql.NullString
	Times_Done     sql.NullFloat64
	Times_Required float64
	Completed_At   sql.NullString

	ty.LongTermTasksCfg
}

type DBYear_Month struct {
	Id    int
	Year  int
	Month int
}
