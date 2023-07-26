create database mydb;
go

use mydb;
go

CREATE TABLE my_table (
    my_nchar NCHAR(10),
    my_char CHAR(10),
    my_nvarchar_max NVARCHAR(MAX),
    my_nvarchar NVARCHAR(50),
    my_varchar_max VARCHAR(MAX),
    my_varchar VARCHAR(50),
    my_ntext NTEXT,
    my_text TEXT,
    my_bigint BIGINT,
    my_int INT,
    my_smallint SMALLINT,
    my_tinyint TINYINT,
    my_float FLOAT,
    my_real REAL,
    my_decimal DECIMAL(10, 5),
    my_money MONEY,
    my_smallmoney SMALLMONEY,
    my_datetime2 DATETIME2,
    my_datetime DATETIME,
    my_smalldatetime SMALLDATETIME,
    my_datetimeoffset DATETIMEOFFSET,
    my_date DATE,
    my_time TIME,
    my_binary BINARY(50),
    my_varbinary VARBINARY(MAX),
    my_image IMAGE,
    my_bit BIT,
    my_uniqueidentifier UNIQUEIDENTIFIER,
    my_xml XML,
    -- my_hierarchyid HIERARCHYID,
    -- my_sql_variant SQL_VARIANT,
    -- my_rowversion ROWVERSION,
    -- my_geometry GEOMETRY,
    -- my_geography GEOGRAPHY
);
go


-- First row with data
INSERT INTO my_table (
    my_nchar, my_char, my_nvarchar_max, my_nvarchar, my_varchar_max, 
    my_varchar, my_ntext, my_text, my_bigint, my_int, my_smallint, 
    my_tinyint, my_float, my_real, my_decimal, my_money, my_smallmoney, 
    my_datetime2, my_datetime, my_smalldatetime, my_datetimeoffset, 
    my_date, my_time, my_binary, my_varbinary, my_image, my_bit, 
    my_uniqueidentifier, my_xml
    -- my_hierarchyid, my_sql_variant, 
    -- my_geometry, my_geography
    ) 
VALUES (
    N'ABCD', 
    'ABCD', 
    N'This is a test message', 
    N'Test message',
    'This is a test message', 
    'Test message',  
    N'Test ntext', 
    'Test text', 
    123456789, 
    12345, 
    123, 
    12, 
    123.45, 
    123.45, 
    123.43, 
    123.45, 
    123.45, 
    '2023-07-23T14:30:00', 
    '2023-07-23T14:30:00', 
    '2023-07-23T14:30:00', 
    '2023-07-23T14:30:00+00:00', 
    '2023-07-23', 
    '14:30:00', 
    0x010101, 
    0x010101, 
    0x010101, 
    1, 
    NEWID(), 
    '<root><test>Some XML data</test></root>'
    -- hierarchyid::GetRoot(), 
    -- 123, 
    -- geometry::STGeomFromText('POINT (100 100)', 4326), 
    -- geography::STGeomFromText('POINT (-50 50)', 4326)
), (
    N'ABCD', 
    'ABCD', 
    N'This is a test message', 
    N'Test message',
    'This is a test message', 
    'Test message',  
    N'Test ntext', 
    'Test text', 
    123456789, 
    12345, 
    123, 
    12, 
    123.45, 
    123.45, 
    123.43, 
    123.45, 
    123.45, 
    '2023-07-23T14:30:00', 
    '2023-07-23T14:30:00', 
    '2023-07-23T14:30:00', 
    '2023-07-23T14:30:00+00:00', 
    '2023-07-23', 
    '14:30:00', 
    0x010101, 
    0x010101, 
    0x010101, 
    1, 
    NEWID(), 
    '<root><test>Some XML data</test></root>'
    -- hierarchyid::GetRoot(), 
    -- 123, 
    -- geometry::STGeomFromText('POINT (100 100)', 4326), 
    -- geography::STGeomFromText('POINT (-50 50)', 4326)
);
go

-- Second row with all NULLs
INSERT INTO my_table (
    my_nchar, my_char, my_nvarchar_max, my_nvarchar, my_varchar_max, 
    my_varchar, my_ntext, my_text, my_bigint, my_int, my_smallint, 
    my_tinyint, my_float, my_real, my_decimal, my_money, my_smallmoney, 
    my_datetime2, my_datetime, my_smalldatetime, my_datetimeoffset, 
    my_date, my_time, my_binary, my_varbinary, my_image, my_bit, 
    my_uniqueidentifier, my_xml
    -- my_hierarchyid, my_sql_variant, 
    -- my_geometry, my_geography
    ) 
VALUES (
    NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 
    NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, 
    NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL
    -- NULL, NULL, NULL, NULL
);
go