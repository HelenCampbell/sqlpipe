package engine

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/calmitchell617/sqlpipe/internal/data"
	_ "github.com/denisenkom/go-mssqldb"
)

var mssql *sql.DB

type MSSQL struct {
	dsType          string
	driverName      string `json:"-"`
	connString      string `json:"-"`
	debugConnString string
	db              *sql.DB
}

func (dsConn MSSQL) execute(query string) (rows *sql.Rows, errProperties map[string]string, err error) {
	return standardExecute(query, dsConn.dsType, dsConn.db)
}

func (dsConn MSSQL) workerPoolExecute(ch chan string, wg *sync.WaitGroup) {
	for query := range ch {
		rows, errProperties, err := dsConn.execute(query)
		if err != nil {
			fmt.Println(err, errProperties)
		}
		rows.Close()
	}
	wg.Done()
}

func (dsConn MSSQL) closeDb() {
	dsConn.db.Close()
}

func (dsConn MSSQL) writeSyncInsert(
	row []string,
	relation relation,
	rowsColumnInfo RowsColumnInfo,
) (
	query string,
) {
	return standardWriteSyncInsert(dsConn, row, relation, rowsColumnInfo)
}

func getNewMSSQL(
	connection data.Connection,
) (
	dsConn DsConnection,
	errProperties map[string]string,
	err error,
) {

	connString := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%v?database=%s",
		connection.Username,
		connection.Password,
		connection.Hostname,
		connection.Port,
		connection.DbName,
	)

	mssql, err = sql.Open("mssql", connString)

	if err != nil {
		return dsConn, errProperties, err
	}

	dsConn = MSSQL{
		"mssql",
		"mssql",
		fmt.Sprintf(
			"sqlserver://%s:%s@%s:%v?database=%s",
			connection.Username,
			connection.Password,
			connection.Hostname,
			connection.Port,
			connection.DbName,
		),
		fmt.Sprintf(
			"sqlserver://<USERNAME_MASKED>:<PASSWORD_MASKED>@%s:%v?database=%s",
			connection.Hostname,
			connection.Port,
			connection.DbName,
		),
		mssql,
	}

	return dsConn, errProperties, err
}

func (dsConn MSSQL) getRows(
	transfer data.Transfer,
) (
	rows *sql.Rows,
	rowColumnInfo RowsColumnInfo,
	errProperties map[string]string,
	err error,
) {
	return standardGetRows(dsConn, transfer)
}

func (dsConn MSSQL) getFormattedResults(
	query string,
) (
	queryResult QueryResult,
	errProperties map[string]string,
	err error,
) {
	return standardGetFormattedResults(dsConn, query)
}

func (dsConn MSSQL) getIntermediateType(
	colTypeFromDriver string,
) (
	intermediateType string,
	errProperties map[string]string,
	err error,
) {
	switch colTypeFromDriver {
	case "BIGINT":
		intermediateType = "MSSQL_BIGINT"
	case "BIT":
		intermediateType = "MSSQL_BIT"
	case "DECIMAL":
		intermediateType = "MSSQL_DECIMAL"
	case "INT":
		intermediateType = "MSSQL_INT"
	case "MONEY":
		intermediateType = "MSSQL_MONEY"
	case "SMALLINT":
		intermediateType = "MSSQL_SMALLINT"
	case "SMALLMONEY":
		intermediateType = "MSSQL_SMALLMONEY"
	case "TINYINT":
		intermediateType = "MSSQL_TINYINT"
	case "FLOAT":
		intermediateType = "MSSQL_FLOAT"
	case "REAL":
		intermediateType = "MSSQL_REAL"
	case "DATE":
		intermediateType = "MSSQL_DATE"
	case "DATETIME2":
		intermediateType = "MSSQL_DATETIME2"
	case "DATETIME":
		intermediateType = "MSSQL_DATETIME"
	case "DATETIMEOFFSET":
		intermediateType = "MSSQL_DATETIMEOFFSET"
	case "SMALLDATETIME":
		intermediateType = "MSSQL_SMALLDATETIME"
	case "TIME":
		intermediateType = "MSSQL_TIME"
	case "CHAR":
		intermediateType = "MSSQL_CHAR"
	case "VARCHAR":
		intermediateType = "MSSQL_VARCHAR"
	case "TEXT":
		intermediateType = "MSSQL_TEXT"
	case "NCHAR":
		intermediateType = "MSSQL_NCHAR"
	case "NVARCHAR":
		intermediateType = "MSSQL_NVARCHAR"
	case "NTEXT":
		intermediateType = "MSSQL_NTEXT"
	case "BINARY":
		intermediateType = "MSSQL_BINARY"
	case "VARBINARY":
		intermediateType = "MSSQL_VARBINARY"
	case "UNIQUEIDENTIFIER":
		intermediateType = "MSSQL_UNIQUEIDENTIFIER"
	case "XML":
		intermediateType = "MSSQL_XML"
	default:
		err = fmt.Errorf("no MSSQL intermediate type for driver type '%v'", colTypeFromDriver)
	}
	return intermediateType, errProperties, err
}

func (dsConn MSSQL) getConnectionInfo() (dsType string, driveName string, connString string) {
	return dsConn.dsType, dsConn.driverName, dsConn.connString
}

func (dsConn MSSQL) GetDebugInfo() (string, string) {
	return dsConn.dsType, dsConn.debugConnString
}

func (dsConn MSSQL) turboTransfer(
	rows *sql.Rows,
	transfer data.Transfer,
	rowColumnInfo RowsColumnInfo,
) (
	errProperties map[string]string,
	err error,
) {
	return errProperties, err
}

func (dsConn MSSQL) turboWriteMidVal(
	valType string,
	value interface{},
	builder *strings.Builder,
) {
}

func (dsConn MSSQL) turboWriteEndVal(
	valType string,
	value interface{},
	builder *strings.Builder,
) {
}

func (db MSSQL) insertChecker(currentLen int, currentRow int) bool {
	if currentRow%1000 == 0 {
		return true
	} else {
		return false
	}
}

func (dsConn MSSQL) dropTable(
	transfer data.Transfer,
) (
	errProperties map[string]string,
	err error,
) {
	return dropTableIfExistsWithSchema(dsConn, transfer)
}

func (dsConn MSSQL) deleteFromTable(
	transfer data.Transfer,
) (
	errProperties map[string]string,
	err error,
) {
	return deleteFromTableWithSchema(dsConn, transfer)
}

func (dsConn MSSQL) createTable(
	transfer data.Transfer,
	columnInfo RowsColumnInfo,
) (
	errProperties map[string]string,
	err error,
) {
	return standardCreateTable(dsConn, transfer, columnInfo)
}

func (dsConn MSSQL) getValToWriteMidRow(valType string, value interface{}) string {
	return mssqlValWriters[valType](value, ",")
}

func (dsConn MSSQL) getValToWriteRaw(valType string, value interface{}) string {
	return mssqlValWriters[valType](value, "")
}

func (dsConn MSSQL) getValToWriteRowEnd(valType string, value interface{}) string {
	return mssqlValWriters[valType](value, ")")
}

func (dsConn MSSQL) getRowStarter() string {
	return standardGetRowStarter()
}

func (dsConn MSSQL) getQueryStarter(targetTable string, columnInfo RowsColumnInfo) string {
	return standardGetQueryStarter(targetTable, columnInfo.ColumnNames)
}

func mssqlWriteBit(value interface{}, terminator string) string {

	var returnVal string

	switch v := value.(type) {
	case bool:
		if v {
			returnVal = fmt.Sprintf("1%s", terminator)
		} else {
			returnVal = fmt.Sprintf("0%s", terminator)
		}
	default:
		return fmt.Sprintf("null%s", terminator)
	}
	return returnVal
}

func mssqlWriteHexBytes(value interface{}, terminator string) string {
	return fmt.Sprintf("CONVERT(VARBINARY(8000), '0x%x', 1)%s", value, terminator)
}

func mssqlWriteUniqueIdentifier(value interface{}, terminator string) string {
	// This is a really stupid fix but it works
	// https://github.com/denisenkom/go-mssqldb/issues/56
	// I guess the bits get shifted around in the first half of these bytes... whatever
	var returnVal string

	switch v := value.(type) {
	case []uint8:
		returnVal = fmt.Sprintf("N'%X%X%X%X-%X%X-%X%X-%X%X-%X'%s",
			v[3],
			v[2],
			v[1],
			v[0],
			v[5],
			v[4],
			v[7],
			v[6],
			v[8],
			v[9],
			v[10:],
			terminator,
		)
	default:
		return fmt.Sprintf("null%s", terminator)
	}
	return returnVal
}

func (dsConn MSSQL) getQueryEnder(targetTable string) string {
	return ""
}

func mssqlWriteDateTime(value interface{}, terminator string) string {
	var returnVal string

	switch v := value.(type) {
	case time.Time:
		returnVal = fmt.Sprintf("CONVERT(DATETIME2, '%s', 121)%s", v.Format("2006-01-02 15:04:05.0000000"), terminator)
	default:
		return fmt.Sprintf("null%s", terminator)
	}

	return returnVal
}

func mssqlWriteDateTimeFromPostgreSQLSync(value interface{}, terminator string) string {
	value = strings.Split(fmt.Sprint(value), "+")[0]
	return fmt.Sprintf("CONVERT(DATETIME2, '%s', 121)%s", value, terminator)
}

func mssqlWriteDateTimeWithTZ(value interface{}, terminator string) string {
	var returnVal string

	switch v := value.(type) {
	case time.Time:
		returnVal = fmt.Sprintf("CONVERT(DATETIMEOFFSET, '%s', 121)%s", v.Format("2006-01-02 15:04:05.000 -07:00"), terminator)
	default:
		return fmt.Sprintf("null%s", terminator)
	}

	return returnVal
}

func mssqlWriteTime(value interface{}, terminator string) string {
	var returnVal string

	switch v := value.(type) {
	case time.Time:
		returnVal = fmt.Sprintf("CONVERT(TIME, '%s', 121)%s", v.Format("15:04:05.000"), terminator)
	default:
		return fmt.Sprintf("null%s", terminator)
	}

	return returnVal
}

func (dsConn MSSQL) getCreateTableType(
	resultSetColInfo RowsColumnInfo,
	colNum int,
) (
	createType string,
) {
	scanType := resultSetColInfo.ColumnScanTypes[colNum]
	intermediateType := resultSetColInfo.ColumnIntermediateTypes[colNum]

	switch scanType.Name() {
	// Generics
	case "bool":
		createType = "BIT"
	case "int", "int8", "int16", "int32", "uint8", "uint16":
		createType = "INT"
	case "int64", "uint32", "uint64":
		createType = "BIGINT"
	case "float32":
		createType = "REAL"
	case "float64":
		createType = "FLOAT"
	case "Time":
		createType = "DATETIME2"
	}

	if createType != "" {
		return createType
	}

	switch intermediateType {
	// PostgreSQL
	case "PostgreSQL_BIGINT":
		createType = "BIGINT"
	case "PostgreSQL_BIT":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_VARBIT":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_BOOLEAN":
		createType = "BIT"
	case "PostgreSQL_BOX":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_BYTEA":
		createType = "VARBINARY(8000)"
	case "PostgreSQL_BPCHAR":
		createType = "NVARCHAR(4000)"
	case "PostgreSQL_CIDR":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_CIRCLE":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_DATE":
		createType = "DATE"
	case "PostgreSQL_FLOAT8":
		createType = "FLOAT"
	case "PostgreSQL_INET":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_INT4":
		createType = "INT"
	case "PostgreSQL_INTERVAL":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_JSON":
		createType = "NVARCHAR(4000)"
	case "PostgreSQL_JSONB":
		createType = "NVARCHAR(4000)"
	case "PostgreSQL_LINE":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_LSEG":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_MACADDR":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_MONEY":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_PATH":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_PG_LSN":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_POINT":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_POLYGON":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_FLOAT4":
		createType = "REAL"
	case "PostgreSQL_INT2":
		createType = "SMALLINT"
	case "PostgreSQL_TEXT":
		createType = "NTEXT"
	case "PostgreSQL_TIME":
		createType = "TIME"
	case "PostgreSQL_TIMETZ":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_TIMESTAMP":
		createType = "DATETIME2"
	case "PostgreSQL_TIMESTAMPTZ":
		createType = "DATETIMEOFFSET"
	case "PostgreSQL_TSQUERY":
		createType = "NVARCHAR(4000)"
	case "PostgreSQL_TSVECTOR":
		createType = "NVARCHAR(4000)"
	case "PostgreSQL_TXID_SNAPSHOT":
		createType = "VARCHAR(8000)"
	case "PostgreSQL_UUID":
		createType = "UNIQUEIDENTIFIER"
	case "PostgreSQL_XML":
		createType = "XML"
	case "PostgreSQL_VARCHAR":
		createType = fmt.Sprintf(
			"NVARCHAR(%d)",
			resultSetColInfo.ColumnLengths[colNum],
		)
	case "PostgreSQL_DECIMAL":
		createType = fmt.Sprintf(
			"DECIMAL(%d,%d)",
			resultSetColInfo.ColumnPrecisions[colNum],
			resultSetColInfo.ColumnScales[colNum],
		)

	// MySQL
	case "MySQL_BIT":
		createType = "VARCHAR(8000)"
	case "MySQL_TINYINT":
		createType = "TINYINT"
	case "MySQL_SMALLINT":
		createType = "SMALLINT"
	case "MySQL_MEDIUMINT":
		createType = "INT"
	case "MySQL_INT":
		createType = "INT"
	case "MySQL_BIGINT":
		createType = "BIGINT"
	case "MySQL_FLOAT4":
		createType = "REAL"
	case "MySQL_FLOAT8":
		createType = "FLOAT"
	case "MySQL_DATE":
		createType = "DATE"
	case "MySQL_TIME":
		createType = "TIME"
	case "MySQL_DATETIME":
		createType = "DATETIME2"
	case "MySQL_TIMESTAMP":
		createType = "DATETIME2"
	case "MySQL_YEAR":
		createType = "INT"
	case "MySQL_CHAR":
		createType = "NVARCHAR(255)"
	case "MySQL_VARCHAR":
		createType = "NVARCHAR(4000)"
	case "MySQL_TEXT":
		createType = "NTEXT"
	case "MySQL_BINARY":
		createType = "VARBINARY(255)"
	case "MySQL_VARBINARY":
		createType = "VARBINARY(8000)"
	case "MySQL_BLOB":
		createType = "VARBINARY(8000)"
	case "MySQL_GEOMETRY":
		createType = "VARBINARY(8000)"
	case "MySQL_JSON":
		createType = "NVARCHAR(4000)"
	case "MySQL_DECIMAL":
		createType = fmt.Sprintf(
			"DECIMAL(%d,%d)",
			resultSetColInfo.ColumnPrecisions[colNum],
			resultSetColInfo.ColumnScales[colNum],
		)
	case "MSSQL_BIGINT":
		createType = "BIGINT"
	case "MSSQL_BIT":
		createType = "BIT"
	case "MSSQL_INT":
		createType = "INT"
	case "MSSQL_MONEY":
		createType = "VARCHAR(8000)"
	case "MSSQL_SMALLINT":
		createType = "SMALLINT"
	case "MSSQL_SMALLMONEY":
		createType = "VARCHAR(8000)"
	case "MSSQL_TINYINT":
		createType = "TINYINT"
	case "MSSQL_FLOAT":
		createType = "FLOAT"
	case "MSSQL_REAL":
		createType = "REAL"
	case "MSSQL_DATE":
		createType = "DATE"
	case "MSSQL_DATETIME2":
		createType = "DATETIME2"
	case "MSSQL_DATETIME":
		createType = "DATETIME"
	case "MSSQL_DATETIMEOFFSET":
		createType = "DATETIMEOFFSET"
	case "MSSQL_SMALLDATETIME":
		createType = "SMALLDATETIME"
	case "MSSQL_TIME":
		createType = "DATETIME"
	case "MSSQL_TEXT":
		createType = "TEXT"
	case "MSSQL_NTEXT":
		createType = "NTEXT"
	case "MSSQL_BINARY":
		createType = "VARBINARY(8000)"
	case "MSSQL_UNIQUEIDENTIFIER":
		createType = "UNIQUEIDENTIFIER"
	case "MSSQL_XML":
		createType = "XML"
	case "MSSQL_DECIMAL":
		createType = fmt.Sprintf(
			"DECIMAL(%d,%d)",
			resultSetColInfo.ColumnPrecisions[colNum],
			resultSetColInfo.ColumnScales[colNum],
		)
	case "MSSQL_CHAR":
		createType = fmt.Sprintf(
			"CHAR(%d)",
			resultSetColInfo.ColumnLengths[colNum],
		)
	case "MSSQL_VARCHAR":
		createType = fmt.Sprintf(
			"VARCHAR(%d)",
			resultSetColInfo.ColumnLengths[colNum],
		)
	case "MSSQL_NCHAR":
		createType = fmt.Sprintf(
			"NCHAR(%d)",
			resultSetColInfo.ColumnLengths[colNum],
		)
	case "MSSQL_NVARCHAR":
		createType = fmt.Sprintf(
			"NVARCHAR(%d)",
			resultSetColInfo.ColumnLengths[colNum],
		)
	case "MSSQL_VARBINARY":
		createType = fmt.Sprintf(
			"VARBINARY(%d)",
			resultSetColInfo.ColumnLengths[colNum],
		)

	case "Oracle_OCIClobLocator":
		createType = "NVARCHAR(4000)"
	case "Oracle_OCIBlobLocator":
		createType = "VARBINARY(8000)"
	case "Oracle_LONG":
		createType = "NTEXT"
	case "Oracle_NUMBER":
		createType = "FLOAT"
	case "Oracle_DATE":
		createType = "DATE"
	case "Oracle_TimeStampDTY":
		createType = "DATETIME2"
	case "Oracle_CHAR":
		createType = "NTEXT"
	case "Oracle_NCHAR":
		createType = "NTEXT"

	// SNOWFLAKE

	case "Snowflake_NUMBER":
		createType = "FLOAT"
	case "Snowflake_BINARY":
		createType = "VARBINARY(8000)"
	case "Snowflake_REAL":
		createType = "FLOAT"
	case "Snowflake_TEXT":
		createType = "NVARCHAR(4000)"
	case "Snowflake_BOOLEAN":
		createType = "BIT"
	case "Snowflake_DATE":
		createType = "DATE"
	case "Snowflake_TIME":
		createType = "TIME"
	case "Snowflake_TIMESTAMP_LTZ":
		createType = "DATETIMEOFFSET"
	case "Snowflake_TIMESTAMP_NTZ":
		createType = "DATETIME2"
	case "Snowflake_TIMESTAMP_TZ":
		createType = "DATETIMEOFFSET"
	case "Snowflake_VARIANT":
		createType = "NVARCHAR(4000)"
	case "Snowflake_OBJECT":
		createType = "NVARCHAR(4000)"
	case "Snowflake_ARRAY":
		createType = "NVARCHAR(4000)"

	// Redshift

	case "Redshift_BIGINT":
		createType = "BIGINT"
	case "Redshift_BOOLEAN":
		createType = "BIT"
	case "Redshift_CHAR":
		createType = "NVARCHAR(4000)"
	case "Redshift_BPCHAR":
		createType = "NVARCHAR(MAX)"
	case "Redshift_DATE":
		createType = "DATE"
	case "Redshift_DOUBLE":
		createType = "FLOAT"
	case "Redshift_INT":
		createType = "INT"
	case "Redshift_REAL":
		createType = "REAL"
	case "Redshift_SMALLINT":
		createType = "SMALLINT"
	case "Redshift_TIME":
		createType = "TIME"
	case "Redshift_TIMETZ":
		createType = "NVARCHAR(4000)"
	case "Redshift_TIMESTAMP":
		createType = "DATETIME2"
	case "Redshift_TIMESTAMPTZ":
		createType = "DATETIMEOFFSET"
	case "Redshift_NUMERIC":
		createType = "FLOAT"
	case "Redshift_VARCHAR":
		createType = fmt.Sprintf(
			"NVARCHAR(%d)",
			resultSetColInfo.ColumnLengths[colNum],
		)
	default:
		createType = "NTEXT"
	}

	return createType
}

var mssqlValWriters = map[string]func(value interface{}, terminator string) string{

	// Generics
	"bool":    mssqlWriteBit,
	"float32": writeInsertFloat,
	"float64": writeInsertFloat,
	"int16":   writeInsertInt,
	"int32":   writeInsertInt,
	"int64":   writeInsertInt,
	"Time":    mssqlWriteDateTime,

	// PostgreSQL
	"PostgreSQL_BIGINT":        writeInsertInt,
	"PostgreSQL_BIT":           writeInsertStringNoEscape,
	"PostgreSQL_VARBIT":        writeInsertStringNoEscape,
	"PostgreSQL_BOOLEAN":       mssqlWriteBit,
	"PostgreSQL_BOX":           writeInsertStringNoEscape,
	"PostgreSQL_BYTEA":         mssqlWriteHexBytes,
	"PostgreSQL_CIDR":          writeInsertStringNoEscape,
	"PostgreSQL_CIRCLE":        writeInsertStringNoEscape,
	"PostgreSQL_FLOAT8":        writeInsertFloat,
	"PostgreSQL_INET":          writeInsertStringNoEscape,
	"PostgreSQL_INT4":          writeInsertInt,
	"PostgreSQL_INTERVAL":      writeInsertStringNoEscape,
	"PostgreSQL_LINE":          writeInsertStringNoEscape,
	"PostgreSQL_LSEG":          writeInsertStringNoEscape,
	"PostgreSQL_MACADDR":       writeInsertStringNoEscape,
	"PostgreSQL_MONEY":         writeInsertStringNoEscape,
	"PostgreSQL_DECIMAL":       writeInsertRawStringNoQuotes,
	"PostgreSQL_PATH":          writeInsertStringNoEscape,
	"PostgreSQL_PG_LSN":        writeInsertStringNoEscape,
	"PostgreSQL_POINT":         writeInsertStringNoEscape,
	"PostgreSQL_POLYGON":       writeInsertStringNoEscape,
	"PostgreSQL_FLOAT4":        writeInsertFloat,
	"PostgreSQL_INT2":          writeInsertInt,
	"PostgreSQL_TIME":          writeInsertStringNoEscape,
	"PostgreSQL_TIMETZ":        writeInsertStringNoEscape,
	"PostgreSQL_TXID_SNAPSHOT": writeInsertStringNoEscape,
	"PostgreSQL_UUID":          writeInsertStringNoEscape,
	"PostgreSQL_VARCHAR":       writeInsertEscapedString,
	"PostgreSQL_BPCHAR":        writeInsertEscapedString,
	"PostgreSQL_DATE":          mssqlWriteDateTime,
	"PostgreSQL_JSON":          writeInsertEscapedString,
	"PostgreSQL_JSONB":         writeInsertEscapedString,
	"PostgreSQL_TEXT":          writeInsertEscapedString,
	"PostgreSQL_TIMESTAMP":     mssqlWriteDateTime,
	"PostgreSQL_TIMESTAMPTZ":   mssqlWriteDateTimeWithTZ,
	"PostgreSQL_TSQUERY":       writeInsertEscapedString,
	"PostgreSQL_TSVECTOR":      writeInsertEscapedString,
	"PostgreSQL_XML":           writeInsertEscapedString,
	// Syncs
	"PostgreSQL_BIGINT_SYNC":      writeInsertRawStringNoQuotes,
	"PostgreSQL_BOOL_SYNC":        writeNumberFromPostgreSQLBoolSync,
	"PostgreSQL_DATE_SYNC":        writeInsertStringNoEscape,
	"PostgreSQL_DOUBLE_SYNC":      writeInsertRawStringNoQuotes,
	"PostgreSQL_INT_SYNC":         writeInsertRawStringNoQuotes,
	"PostgreSQL_FLOAT_SYNC":       writeInsertRawStringNoQuotes,
	"PostgreSQL_SMALLINT_SYNC":    writeInsertRawStringNoQuotes,
	"PostgreSQL_TIMESTAMP_SYNC":   mssqlWriteDateTimeFromPostgreSQLSync,
	"PostgreSQL_TIMESTAMPTZ_SYNC": mssqlWriteDateTimeFromPostgreSQLSync,
	"NIL":                         postgresqlWriteNone,

	// MYSQL

	"MySQL_BIT":       writeInsertBinaryString,
	"MySQL_TINYINT":   writeInsertRawStringNoQuotes,
	"MySQL_SMALLINT":  writeInsertRawStringNoQuotes,
	"MySQL_MEDIUMINT": writeInsertRawStringNoQuotes,
	"MySQL_INT":       writeInsertRawStringNoQuotes,
	"MySQL_BIGINT":    writeInsertRawStringNoQuotes,
	"MySQL_DECIMAL":   writeInsertRawStringNoQuotes,
	"MySQL_FLOAT4":    writeInsertRawStringNoQuotes,
	"MySQL_FLOAT8":    writeInsertRawStringNoQuotes,
	"MySQL_DATE":      writeInsertStringNoEscape,
	"MySQL_TIME":      writeInsertStringNoEscape,
	"MySQL_TIMESTAMP": writeInsertStringNoEscape,
	"MySQL_DATETIME":  writeInsertStringNoEscape,
	"MySQL_YEAR":      writeInsertRawStringNoQuotes,
	"MySQL_CHAR":      writeInsertEscapedString,
	"MySQL_VARCHAR":   writeInsertEscapedString,
	"MySQL_TEXT":      writeInsertEscapedString,
	"MySQL_BINARY":    mssqlWriteHexBytes,
	"MySQL_VARBINARY": mssqlWriteHexBytes,
	"MySQL_BLOB":      mssqlWriteHexBytes,
	"MySQL_GEOMETRY":  mssqlWriteHexBytes,
	"MySQL_JSON":      writeInsertEscapedString,

	// MSSQL

	"MSSQL_BIGINT":           writeInsertInt,
	"MSSQL_BIT":              mssqlWriteBit,
	"MSSQL_DECIMAL":          writeInsertRawStringNoQuotes,
	"MSSQL_INT":              writeInsertInt,
	"MSSQL_MONEY":            writeInsertStringNoEscape,
	"MSSQL_SMALLINT":         writeInsertInt,
	"MSSQL_SMALLMONEY":       writeInsertStringNoEscape,
	"MSSQL_TINYINT":          writeInsertInt,
	"MSSQL_FLOAT":            writeInsertFloat,
	"MSSQL_REAL":             writeInsertFloat,
	"MSSQL_DATE":             mssqlWriteDateTime,
	"MSSQL_DATETIME2":        mssqlWriteDateTime,
	"MSSQL_DATETIME":         mssqlWriteDateTime,
	"MSSQL_DATETIMEOFFSET":   mssqlWriteDateTime,
	"MSSQL_SMALLDATETIME":    mssqlWriteDateTime,
	"MSSQL_TIME":             mssqlWriteDateTime,
	"MSSQL_CHAR":             writeInsertEscapedString,
	"MSSQL_VARCHAR":          writeInsertEscapedString,
	"MSSQL_TEXT":             writeInsertEscapedString,
	"MSSQL_NCHAR":            writeInsertEscapedString,
	"MSSQL_NVARCHAR":         writeInsertEscapedString,
	"MSSQL_NTEXT":            writeInsertEscapedString,
	"MSSQL_BINARY":           mssqlWriteHexBytes,
	"MSSQL_VARBINARY":        mssqlWriteHexBytes,
	"MSSQL_UNIQUEIDENTIFIER": mssqlWriteUniqueIdentifier,
	"MSSQL_XML":              writeInsertEscapedString,

	// Oracle

	"Oracle_CHAR":           writeInsertEscapedString,
	"Oracle_NCHAR":          writeInsertEscapedString,
	"Oracle_OCIClobLocator": writeInsertEscapedString,
	"Oracle_OCIBlobLocator": mssqlWriteHexBytes,
	"Oracle_LONG":           writeInsertEscapedString,
	"Oracle_NUMBER":         oracleWriteNumber,
	"Oracle_DATE":           mssqlWriteDateTime,
	"Oracle_TimeStampDTY":   mssqlWriteDateTime,

	// Snowflake

	"Snowflake_NUMBER":        writeInsertRawStringNoQuotes,
	"Snowflake_REAL":          writeInsertRawStringNoQuotes,
	"Snowflake_TEXT":          writeInsertEscapedString,
	"Snowflake_BOOLEAN":       writeInsertStringNoEscape,
	"Snowflake_DATE":          mssqlWriteDateTime,
	"Snowflake_TIME":          mssqlWriteTime,
	"Snowflake_TIMESTAMP_LTZ": mssqlWriteDateTime,
	"Snowflake_TIMESTAMP_NTZ": mssqlWriteDateTime,
	"Snowflake_TIMESTAMP_TZ":  mssqlWriteDateTime,
	"Snowflake_VARIANT":       writeInsertEscapedStringRemoveNewines,
	"Snowflake_OBJECT":        writeInsertEscapedStringRemoveNewines,
	"Snowflake_ARRAY":         writeInsertEscapedStringRemoveNewines,
	"Snowflake_BINARY":        mssqlWriteHexBytes,

	// Redshift

	"Redshift_BIGINT":      writeInsertInt,
	"Redshift_BOOLEAN":     mssqlWriteBit,
	"Redshift_CHAR":        writeInsertEscapedString,
	"Redshift_VARCHAR":     writeInsertEscapedString,
	"Redshift_BPCHAR":      writeInsertEscapedString,
	"Redshift_DATE":        mssqlWriteDateTime,
	"Redshift_DOUBLE":      writeInsertFloat,
	"Redshift_INT":         writeInsertInt,
	"Redshift_NUMERIC":     writeInsertRawStringNoQuotes,
	"Redshift_REAL":        writeInsertFloat,
	"Redshift_SMALLINT":    writeInsertInt,
	"Redshift_TIME":        writeInsertStringNoEscape,
	"Redshift_TIMETZ":      writeInsertStringNoEscape,
	"Redshift_TIMESTAMP":   mssqlWriteDateTime,
	"Redshift_TIMESTAMPTZ": mssqlWriteDateTime,
}
