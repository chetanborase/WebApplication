package database

import (
	"WebApplication/config"
)

var (
	totalRows = 0
	query     = `select count(UserInfo.id) from UserInfo where DeletedAt IS NULL`
)

//validating maxId PageNo AND PageSize(limit)
func GetTotalRecords() (int, error) {
	row := config.GetDB().QueryRow(query)
	err := row.Scan(&totalRows)
	if err != nil {
		return 0, err
	}
	return totalRows, nil
}
