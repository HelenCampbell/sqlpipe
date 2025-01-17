package formatters

import (
	"database/sql"

	"github.com/sqlpipe/sqlpipe/internal/engine/transfers/formatters/shared"
)

var MssqlCreateFormatters = map[string]func(column *sql.ColumnType, terminator string) (string, error){
	"SQL_UNKNOWN_TYPE":    shared.NTextCreateFormatter,
	"SQL_CHAR":            shared.NTextCreateFormatter,
	"SQL_NUMERIC":         shared.DecimalCreateFormatter,
	"SQL_DECIMAL":         shared.DecimalCreateFormatter,
	"SQL_INTEGER":         shared.IntCreateFormatter,
	"SQL_SMALLINT":        shared.SmallIntCreateFormatter,
	"SQL_FLOAT":           shared.FloatCreateFormatter,
	"SQL_REAL":            shared.FloatCreateFormatter,
	"SQL_DOUBLE":          shared.FloatCreateFormatter,
	"SQL_DATETIME":        shared.Datetime2CreateFormatter,
	"SQL_TIME":            shared.TimeCreateFormatter,
	"SQL_VARCHAR":         shared.NTextCreateFormatter,
	"SQL_TYPE_DATE":       shared.DateCreateFormatter,
	"SQL_TYPE_TIME":       shared.TimeCreateFormatter,
	"SQL_TYPE_TIMESTAMP":  shared.Datetime2CreateFormatter,
	"SQL_TIMESTAMP":       shared.Datetime2CreateFormatter,
	"SQL_LONGVARCHAR":     shared.NTextCreateFormatter,
	"SQL_BINARY":          shared.NTextCreateFormatter,
	"SQL_VARBINARY":       shared.NTextCreateFormatter,
	"SQL_LONGVARBINARY":   shared.NTextCreateFormatter,
	"SQL_BIGINT":          shared.BigIntCreateFormatter,
	"SQL_TINYINT":         shared.SmallIntCreateFormatter,
	"SQL_BIT":             shared.BitCreateFormatter,
	"SQL_WCHAR":           shared.NTextCreateFormatter,
	"SQL_WVARCHAR":        shared.NTextCreateFormatter,
	"SQL_WLONGVARCHAR":    shared.NTextCreateFormatter,
	"SQL_GUID":            shared.UniqueIdentifierCreateFormatter,
	"SQL_SIGNED_OFFSET":   shared.NTextCreateFormatter,
	"SQL_UNSIGNED_OFFSET": shared.NTextCreateFormatter,
	"SQL_SS_XML":          shared.XmlCreateFormatter,
	"SQL_SS_TIME2":        shared.TimeCreateFormatter,
}

var MssqlValFormatters = map[string]func(value interface{}, terminator string) (formattedValue string, err error){
	"SQL_UNKNOWN_TYPE":    shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_CHAR":            shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_NUMERIC":         shared.RawXnull,
	"SQL_DECIMAL":         shared.RawXnull,
	"SQL_INTEGER":         shared.RawXnull,
	"SQL_SMALLINT":        shared.RawXnull,
	"SQL_FLOAT":           shared.RawXnull,
	"SQL_REAL":            shared.RawXnull,
	"SQL_DOUBLE":          shared.RawXnull,
	"SQL_DATETIME":        shared.CastToTimeFormatToTimetampStringXnull,
	"SQL_TIME":            shared.CastToTimeFormatToTimeStringXnull,
	"SQL_VARCHAR":         shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_TYPE_DATE":       shared.CastToTimeFormatToDateStringXnull,
	"SQL_TYPE_TIME":       shared.CastToTimeFormatToTimeStringXnull,
	"SQL_TYPE_TIMESTAMP":  shared.CastToTimeFormatToTimetampStringXnull,
	"SQL_TIMESTAMP":       shared.CastToTimeFormatToTimetampStringXnull,
	"SQL_LONGVARCHAR":     shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_BINARY":          shared.CastToBytesCastToStringPrintQuotedHexXnull,
	"SQL_VARBINARY":       shared.CastToBytesCastToStringPrintQuotedHexXnull,
	"SQL_LONGVARBINARY":   shared.CastToBytesCastToStringPrintQuotedHexXnull,
	"SQL_BIGINT":          shared.RawXnull,
	"SQL_TINYINT":         shared.RawXnull,
	"SQL_BIT":             shared.CastToBoolWriteBinaryEquivalentXnull,
	"SQL_WCHAR":           shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_WVARCHAR":        shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_WLONGVARCHAR":    shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_GUID":            shared.QuotedXnull,
	"SQL_SIGNED_OFFSET":   shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_UNSIGNED_OFFSET": shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_SS_XML":          shared.CastToBytesCastToStringPrintQuotedXnull,
	"SQL_SS_TIME2":        shared.CastToBytesCastToStringPrintQuotedXnull,
}
