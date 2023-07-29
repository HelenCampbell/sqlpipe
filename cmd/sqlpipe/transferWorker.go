package main

import (
	"fmt"
)

type ColumnInfo struct {
	name       string
	pipeType   string
	scanType   string
	decimalOk  bool
	precision  int64
	scale      int64
	lengthOk   bool
	length     int64
	nullableOk bool
	nullable   bool
}

func runTransfer(transfer Transfer) {
	if transfer.DropTargetTable {
		err := transfer.Target.dropTable(transfer.TargetSchema, transfer.TargetTable)
		if err != nil {
			transferError(transfer, fmt.Errorf("error dropping target table :: %v", err))
			return
		}
	}

	var err error
	transfer.rows, err = transfer.Source.query(transfer.Query)
	if err != nil {
		transferError(transfer, fmt.Errorf("error querying source :: %v", err))
		return
	}
	defer transfer.rows.Close()

	transfer.ColumnInfo, err = transfer.Source.getColumnInfo(transfer.rows)
	if err != nil {
		transferError(transfer, fmt.Errorf("error getting source column info :: %v", err))
		return
	}

	if transfer.CreateTargetTable {
		err = transfer.Target.createTable(transfer.TargetSchema, transfer.TargetTable, transfer.ColumnInfo)
		if err != nil {
			transferError(transfer, fmt.Errorf("error creating target table :: %v", err))
			return
		}
	}

	tmpDir, err := transfer.Source.createPipeFiles(transfer.rows, transfer.ColumnInfo, transfer.Id)
	if err != nil {
		transferError(transfer, fmt.Errorf("error writing pipe file :: %v", err))
		return
	}

	err = transfer.Target.insertPipeFiles(tmpDir, transfer.Id, transfer.ColumnInfo, transfer.TargetTable, transfer.TargetSchema)
	if err != nil {
		transferError(transfer, fmt.Errorf("error inserting data :: %v", err))
	}
}
