package db

import (
	"fmt"
	"go.uber.org/zap"
	"strings"
	"git.tianrang-inc.com/data-brain/trains/log"
)

const BatchLoadAmount = 100

// 批量 insert 到 postgres
// updateFields 不为空时，conflict 时可触发更新；否则不更新
func load(tableName string, fields, conflictFields, updateFields []string, records [][]interface{}) (int, error) {

	cursor, err := PgConn.Connect()
	if err != nil {
		return 0, err
	}
	defer PgConn.Close(cursor)

	valueNames := strings.Join(fields, ", ")
	args := make([]interface{}, 1)

	var valuePlaceHolder = strings.Repeat("?,", len(fields))
	valuePlaceHolder = "(" + valuePlaceHolder[:len(valuePlaceHolder)-1] + "),"
	valuePlaceHolders := strings.Repeat(valuePlaceHolder, len(records))
	valuePlaceHolders = valuePlaceHolders[:len(valuePlaceHolders)-1]
	for _, record := range records {
		args = append(args, record...)
	}

	sql := "insert into " + tableName + " (" + valueNames + ") values" + valuePlaceHolders
	if len(conflictFields) > 0 {
		onDups := make([]string, 0)
		sql += " on conflict(" + strings.Join(conflictFields, ", ") + ") do "
		if len(updateFields) > 0 {
			for _, field := range updateFields {
				onDups = append(onDups, field+"=excluded."+field)
			}
			sql += " update set " + strings.Join(onDups, ", ")
		} else {
			sql += " nothing"
		}
	}

	args[0] = sql
	res, err := cursor.Exec(args...)
	log.Logger.Debug(fmt.Sprintf("DEBUG %+v", args))
	if err != nil {
		log.Logger.Error(err.Error())
		return 0, err
	}

	var rowsAffected int64
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		log.Logger.Error("Postgres", zap.Error(err))
	} else if rowsAffected != int64(len(records)) {
		log.Logger.Info("Postgres", zap.Int("to insert", len(records)),
			zap.Int64("inserted", rowsAffected))
	}
	return len(records), nil
}

// 按 BatchLoadAmount 的值批量导入 postgres
func Load(tableName string, fields, conflictFields, updateFields []string, records [][]interface{}) (int, error) {
	bulkSize := BatchLoadAmount
	index := 0
	success := 0
	var count int
	var end int
	var err error
	for index < len(records) {
		end = index + int(bulkSize)
		if end > len(records) {
			end = len(records)
		}
		count, err = load(tableName, fields, conflictFields, updateFields, records[index:end])
		success += count
		index = end
	}
	return success, err
}
