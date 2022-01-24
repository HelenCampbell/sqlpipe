package engine

import (
	"reflect"
	"testing"

	"github.com/calmitchell617/sqlpipe/internal/data"
)

var queryTests = []queryTest{
	// PostgreSQL Setup
	{
		name:                  "postgresqlWideTableDrop",
		connection:            postgresqlTestConnection,
		testQuery:             "drop table if exists wide_table;",
		checkQuery:            "select * from wide_table",
		expectedErr:           "db.Query() threw an error",
		expectedErrProperties: postgresqlWideTableDropErrProperties,
	},
	{
		name:        "postgresqlWideTableCreate",
		connection:  postgresqlTestConnection,
		testQuery:   postgresqlWideTableCreateQuery,
		checkQuery:  "select * from wide_table",
		checkResult: postgresqlWideTableCreateResult,
	},
	{
		name:        "postgresqlWideTableInsert",
		connection:  postgresqlTestConnection,
		testQuery:   postgresqlWideTableInsertQuery,
		checkQuery:  "select * from wide_table",
		checkResult: postgresqlWideTableInsertResult,
	},
	// MSSQL setup
	{
		name:       "mssqlTestingDbDrop",
		connection: mssqlMasterTestConnection,
		testQuery:  "drop database if exists testing",
	},
	{
		name:       "mssqlTestingDbCreate",
		connection: mssqlMasterTestConnection,
		testQuery:  "create database testing",
	},
}

var transferTests = []transferTest{
	// PostgreSQL Transfers
	{
		name:          "postgresql2postgresql_wide",
		source:        postgresqlTestConnection,
		target:        postgresqlTestConnection,
		overwrite:     true,
		targetSchema:  "public",
		targetTable:   "postgresql_wide_table",
		transferQuery: "select * from wide_table",
		checkQuery:    "select * from postgresql_wide_table",
		checkResult:   postgresql2postgresql_wide_result,
	},
	{
		name:          "postgresql2mysql_wide",
		source:        postgresqlTestConnection,
		target:        mysqlTestConnection,
		overwrite:     true,
		targetTable:   "postgresql_wide_table",
		transferQuery: "select * from wide_table",
		checkQuery:    "select * from postgresql_wide_table",
		checkResult:   postgresql2mysql_wide_result,
	},
	{
		name:          "postgresql2mssql_wide",
		source:        postgresqlTestConnection,
		target:        mssqlTestConnection,
		overwrite:     true,
		targetTable:   "postgresql_wide_table",
		transferQuery: "select * from wide_table",
		checkQuery:    "select * from postgresql_wide_table",
		checkResult:   postgresql2mssql_wide_result,
	},
	{
		name:          "postgresql2oracle_wide",
		source:        postgresqlTestConnection,
		target:        oracleTestConnection,
		overwrite:     true,
		targetTable:   "postgresql_wide_table",
		transferQuery: "select * from wide_table",
		checkQuery:    "select * from postgresql_wide_table",
		checkResult:   postgresql2oracle_wide_result,
	},
	{
		name:          "postgresql2redshift_wide",
		source:        postgresqlTestConnection,
		target:        redshiftTestConnection,
		overwrite:     true,
		targetSchema:  "public",
		targetTable:   "postgresql_wide_table",
		transferQuery: "select * from wide_table",
		checkQuery:    "select * from postgresql_wide_table",
		checkResult:   postgresql2redshift_wide_result,
	},
}

func TestRunQuery(t *testing.T) {

	// Loop over the test cases.
	for _, tt := range queryTests {
		t.Run(tt.name, func(t *testing.T) {
			errProperties, err := RunQuery(
				&data.Query{
					Query:      tt.testQuery,
					Connection: tt.connection,
				},
			)

			if err != nil && (err.Error() != tt.expectedErr || !reflect.DeepEqual(errProperties, tt.expectedErrProperties)) {

				t.Fatalf("unable to run test query. err:\n\n%v\n\nerrProperties:\n%#v", err, errProperties)
			}

			if tt.checkQuery != "" {
				dsConn, _, err := GetDs(tt.connection)
				if err != nil {
					t.Fatalf("Couldn't get DsConn")
				}
				queryResult, errProperties, err := standardGetFormattedResults(dsConn, tt.checkQuery)

				if err != nil && err.Error() != tt.expectedErr {
					t.Fatalf("\nwanted error:\n%#v\n\ngot error:\n%#v\n", tt.expectedErr, err)
				}

				if err != nil && !reflect.DeepEqual(errProperties, tt.expectedErrProperties) {
					t.Fatalf("\nwanted errProperties:\n%#v\n\ngot:\n%#v", tt.expectedErrProperties, errProperties)
				}

				if !reflect.DeepEqual(queryResult, tt.checkResult) {
					t.Fatalf("\n\nWanted:\n%#v\n\nGot:\n%#v", tt.checkResult, queryResult)
				}
			}
		})
	}
}

func TestRunTransfer(t *testing.T) {

	// Loop over the test cases.
	for _, tt := range transferTests {

		t.Run(tt.name, func(t *testing.T) {
			errProperties, err := RunTransfer(
				&data.Transfer{
					Query:        tt.transferQuery,
					Overwrite:    tt.overwrite,
					TargetSchema: tt.targetSchema,
					TargetTable:  tt.targetTable,
					Source:       tt.source,
					Target:       tt.target,
				},
			)

			if err != nil {
				t.Fatalf("unable to run transfer. err:\n\n%v\n\nerrProperties:\n%v", err, errProperties)
			}

			if tt.checkQuery != "" {
				dsConn, _, err := GetDs(tt.target)
				if err != nil {
					t.Fatalf("Couldn't get DsConn")
				}
				queryResult, errProperties, err := standardGetFormattedResults(dsConn, tt.checkQuery)

				if err != nil && err.Error() != tt.expectedErr {
					// t.Error(errProperties)
					t.Fatalf("\nwanted error:\n%#v\n\ngot error:\n%#v\n", tt.expectedErr, err)
				}

				if err != nil && !reflect.DeepEqual(errProperties, tt.expectedErrProperties) {
					t.Fatalf("\nwanted errProperties:\n%#v\n\ngot:\n%#v", tt.expectedErrProperties, errProperties)
				}

				if !reflect.DeepEqual(queryResult, tt.checkResult) {
					t.Fatalf("\n\nWanted:\n%#v\n\nGot:\n%#v", tt.checkResult, queryResult)
				}
			}
		})
	}
}
