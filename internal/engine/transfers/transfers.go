package transfers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/sqlpipe/sqlpipe/internal/data"
	"github.com/sqlpipe/sqlpipe/internal/engine/transfers/formatters"
)

func RunTransfer(
	ctx context.Context,
	transfer data.Transfer,
) (
	err error,
) {
	rows, err := transfer.Source.Db.QueryContext(ctx, transfer.Query)
	if err != nil {
		return fmt.Errorf("error running query on source: %v", err.Error())
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error getting column names: %v", err.Error())
	}

	numCols := len(columnNames)

	vals := make([]interface{}, numCols)
	valPtrs := make([]interface{}, numCols)

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return fmt.Errorf("error getting column types: %v", err.Error())
	}
	colDbTypes := []string{}
	for _, colType := range colTypes {
		colDbTypes = append(colDbTypes, colType.DatabaseTypeName())
	}

	schemaSpecifier := ""
	switch transfer.Target.Schema {
	case "":
	default:
		schemaSpecifier = fmt.Sprintf("%v.", transfer.Target.Schema)
	}

	createFormatters := systemCreateFormatters[transfer.Target.SystemType]

	valFormatters := systemValFormatters[transfer.Target.SystemType]

	if transfer.DropTargetTable {
		dropTableCommand := fmt.Sprintf(
			"%v %v%v",
			dropTableCommandStarters[transfer.Target.SystemType],
			schemaSpecifier,
			transfer.Target.Table,
		)

		_, err = transfer.Target.Db.ExecContext(ctx, dropTableCommand)
		if err != nil {
			return fmt.Errorf("error running drop table command: %v", err)
		}
	}

	if transfer.CreateTargetTable {
		createQuery := fmt.Sprintf("create table %v%v(", schemaSpecifier, transfer.Target.Table)

		for i := 0; i < numCols-1; i++ {
			columnSpecifier, err := createFormatters[colDbTypes[i]](colTypes[i], ",")
			if err != nil {
				return fmt.Errorf("error running %v formatter on value %v: %v", colDbTypes[i], colTypes[i], err)
			}
			createQuery = createQuery + columnSpecifier
		}
		columnSpecifier, err := createFormatters[colDbTypes[numCols-1]](colTypes[numCols-1], ")")
		if err != nil {
			return fmt.Errorf("error running %v formatter on value %v: %v", colDbTypes[numCols-1], colTypes[numCols-1], err)
		}
		createQuery = createQuery + columnSpecifier

		_, err = transfer.Target.Db.ExecContext(ctx, createQuery)
		if err != nil {
			return fmt.Errorf("error running create table command: %v", err)
		}
	}

	var batchBuilder strings.Builder

	for i := 0; i < numCols; i++ {
		valPtrs[i] = &vals[i]
	}

	columnNamesString := strings.Join(columnNames, ",")

	isFirstRow := true
	dataRemaining := false
	insertCheckType := insertCheckerTypes[transfer.Target.SystemType]
	insertCheckNum := insertCheckerNums[transfer.Target.SystemType]

	for i := 1; rows.Next(); i++ {
		dataRemaining = true
		rows.Scan(valPtrs...)

		if isFirstRow {
			batchBuilder.WriteString(fmt.Sprintf("insert into %v%v (%v) values (", schemaSpecifier, transfer.Target.Table, columnNamesString))
		} else {
			batchBuilder.WriteString(",(")
		}
		isFirstRow = false
		for j := 0; j < numCols-1; j++ {
			valToWrite, err := valFormatters[colDbTypes[j]](vals[j], ",")
			if err != nil {
				return fmt.Errorf("error running %v formatter on mid-row value %v: %v", colDbTypes[j], vals[j], err)
			}
			batchBuilder.WriteString(valToWrite)
		}
		valToWrite, err := valFormatters[colDbTypes[numCols-1]](vals[numCols-1], ")")
		if err != nil {
			return fmt.Errorf("error running %v formatter on row-end value %v: %v", colDbTypes[numCols-1], vals[numCols-1], err)
		}
		batchBuilder.WriteString(valToWrite)

		switch insertCheckType {
		case "rows":
			if i%insertCheckNum == 0 {
				_, err := transfer.Target.Db.ExecContext(ctx, batchBuilder.String())
				if err != nil {
					return fmt.Errorf("error running mid-batch insert statement: %v", err)
				}

				batchBuilder.Reset()
				isFirstRow = true
				dataRemaining = false
			}
		default:
			if batchBuilder.Len()%insertCheckNum == 0 {
				_, err := transfer.Target.Db.ExecContext(ctx, batchBuilder.String())
				if err != nil {
					return fmt.Errorf("error running mid-batch insert statement: %v", err)
				}

				batchBuilder.Reset()
				isFirstRow = true
				dataRemaining = false
			}
		}

	}

	if dataRemaining {
		stringToWrite := batchBuilder.String()

		_, err := transfer.Target.Db.ExecContext(ctx, stringToWrite)
		if err != nil {
			return fmt.Errorf("error running batch-end insert statement: %v", err)
		}

		batchBuilder.Reset()
	}

	return nil
}

var (
	systemCreateFormatters = map[string]map[string]func(column *sql.ColumnType, terminator string) (string, error){
		"postgresql": formatters.PostgresqlCreateFormatters,
		"mssql":      formatters.MssqlCreateFormatters,
		"mysql":      formatters.MysqlCreateFormatters,
		"snowflake":  formatters.SnowflakeCreateFormatters,
	}
	systemValFormatters = map[string]map[string]func(value interface{}, terminator string) (string, error){
		"postgresql": formatters.PostgresqlValFormatters,
		"mssql":      formatters.MssqlValFormatters,
		"mysql":      formatters.MysqlValFormatters,
		"snowflake":  formatters.SnowflakeValFormatters,
	}
	dropTableCommandStarters = map[string]string{
		"postgresql": "drop table if exists",
		"mssql":      "drop table if exists",
		"mysql":      "drop table if exists",
		"snowflake":  "drop table if exists",
	}
	insertCheckerTypes = map[string]string{
		"postgresql": "length",
		"mysql":      "length",
		"mssql":      "rows",
		"snowflake":  "rows",
	}
	insertCheckerNums = map[string]int{
		"postgresql": 10000000,
		"mysql":      4000000,
		"mssql":      1000,
		"snowflake":  3000,
	}
)
