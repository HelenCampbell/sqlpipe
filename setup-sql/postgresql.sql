create database mydb;

\c mydb;

drop table if exists mydb.public.my_table;

CREATE TABLE mydb.public.my_table (
    my_char CHAR,
    my_varchar VARCHAR,
    my_text TEXT,
    my_int8 INT8,
    my_int4 INT4,
    my_int2 INT2,
    my_float8 FLOAT8,
    my_float4 FLOAT4,
    my_numeric NUMERIC(10, 5),
    my_timestamp TIMESTAMP,
    my_timestamptz TIMESTAMPTZ,
    my_date DATE,
    my_time TIME,
    my_interval INTERVAL,
    my_bytea BYTEA,
    my_box BOX,
    my_circle CIRCLE,
    my_line LINE,
    my_path PATH,
    my_point POINT,
    my_polygon POLYGON,
    my_lseg LSEG,
    my_bool BOOL,
    my_bit BIT,
    my_varbit VARBIT,
    my_uuid UUID,
    my_json JSON,
    my_jsonb JSONB,
    my_inet INET,
    my_macaddr MACADDR,
    my_cidr CIDR,
    my_xml xml,
    my_timetz TIMETZ,
    my_pg_lsn pg_lsn,
    my_tsquery tsquery,
    my_tsvector tsvector,
    my_smallserial SMALLSERIAL,
    my_serial SERIAL,
    my_bigserial BIGSERIAL
);

INSERT INTO mydb.public.my_table (
    my_char,
    my_varchar,
    my_text,
    my_int8,
    my_int4,
    my_int2,
    my_float8,
    my_float4,
    my_numeric,
    my_timestamp,
    my_timestamptz,
    my_date,
    my_time,
    my_interval,
    my_bytea,
    my_box,
    my_circle,
    my_line,
    my_path,
    my_point,
    my_polygon,
    my_lseg,
    my_bool,
    my_bit,
    my_varbit,
    my_uuid,
    my_json,
    my_jsonb,
    my_inet,
    my_macaddr,
    my_cidr,
    my_xml,
    my_timetz,
    my_tsquery,
    my_tsvector,
    my_pg_lsn
)
VALUES (
    'A',
    '"Hello" there',
    'This is ''a tab test    text',
    1234567890123456,
    12345678,
    1234,
    12345678.12345678,
    1234.1234,
    123.456,
    '2023-07-23 12:34:56.643381',
    '2023-07-23 12:34:56.643381+00',
    '2023-07-23',
    '12:34:56.643381',
    '1 day',
    E'\\xDEADBEEF',
    '((1,2),(3,4))',
    '<(1,2),3>',
    '{1,-2,3}',
    '((1,2),(3,4))',
    '(1,2)',
    '((1,2),(3,4))',
    '((1,2),(3,4))',
    true,
    B'1',
    B'101',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    '{"key":"value"}',
    '{"key":"value"}',
    '192.168.1.1',
    '08:00:2B:01:02:03',
    '192.168.1.0/24',
    '<root><test>Some Content</test></root>',
    '12:34:56+00',
    'super & rat',
    'super',
    '1/3B9ACA00'
);

INSERT INTO mydb.public.my_table DEFAULT VALUES;