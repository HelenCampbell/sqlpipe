create database mydb;

\c mydb;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION rdm_ascii(num_chars integer) RETURNS text AS $$
DECLARE
    result text := '';
    i integer := 0;
    random_char integer;
BEGIN
    num_chars = random() * num_chars;

    FOR i IN 1..num_chars LOOP
        random_char := 32 + floor(random() * 95)::integer;
        result := result || chr(random_char);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;


CREATE OR REPLACE FUNCTION rdm_unicode(num_chars integer) RETURNS text AS $$
DECLARE
    result text[]; -- use an array instead of text
    i integer := 0;
    random_char integer;
    range_choice integer;
BEGIN
    num_chars = ceil(random() * num_chars);

    FOR i IN 1..num_chars LOOP
        range_choice := floor(random() * 100)::integer; 

        CASE
            WHEN range_choice < 3 THEN
                random_char := 128 + floor(random() * 128)::integer; -- Extended unicode range
            WHEN range_choice < 4 THEN
                random_char := 256 + floor(random() * 128)::integer; -- Another extended unicode range
            WHEN range_choice < 5 THEN
                random_char := 384 + floor(random() * 112)::integer; -- Another extended unicode range

            WHEN range_choice < 10 THEN
                random_char := 1024 + floor(random() * 256)::integer; -- Cyrillic


            WHEN range_choice < 20 THEN
                random_char := 19968 + floor(random() * 20992)::integer; -- Chinese Han characters

            WHEN range_choice < 25 THEN
                random_char := 1536 + floor(random() * 256)::integer; -- Arabic

            WHEN range_choice < 30 THEN
                random_char := 44032 + floor(random() * 11172)::integer; -- Korean Hangul Syllables


            WHEN range_choice < 32 THEN
                random_char := 12352 + floor(random() * 96)::integer; -- Hiragana
            WHEN range_choice < 34 THEN
                random_char := 12448 + floor(random() * 96)::integer; -- Katakana


            WHEN range_choice < 35 THEN
                random_char := 880 + floor(random() * 128)::integer; -- Greek
            WHEN range_choice < 36 THEN
                random_char := 1424 + floor(random() * 128)::integer; -- Hebrew
            WHEN range_choice < 38 THEN
                random_char := 2304 + floor(random() * 128)::integer; -- Hindi (Devanagari script)
            WHEN range_choice < 39 THEN
                random_char := 3584 + floor(random() * 128)::integer; -- Thai
            WHEN range_choice < 40 THEN
                random_char := 8704 + floor(random() * 256)::integer; -- Mathematical Symbols Basic
            WHEN range_choice < 41 THEN
                random_char := 8192 + floor(random() * 112)::integer; -- General Punctuation
            WHEN range_choice < 42 THEN
                random_char := 9472 + floor(random() * 128)::integer; -- Box Drawing
            WHEN range_choice < 43 THEN
                random_char := 9728 + floor(random() * 192)::integer; -- Dingbats
            WHEN range_choice < 44 THEN
                random_char := 12288 + floor(random() * 64)::integer; -- CJK Symbols and Punctuation
            WHEN range_choice < 45 THEN
                random_char := 7680 + floor(random() * 256)::integer; -- unicode Extended Additional
            WHEN range_choice < 46 THEN
                random_char := 127872 + floor(random() * 256)::integer; -- Symbols for Legacy Computing
            WHEN range_choice < 50 THEN
                random_char := 128512 + floor(random() * 207)::integer; -- Extended mojis
            ELSE -- everything else should be basic unicode range
                random_char := 32 + floor(random() * 95)::integer; -- Basic unicode range
        END CASE;

        result[i] := chr(random_char); -- assign character to the array at position i
    END LOOP;

    RETURN array_to_string(result, ''); -- convert the array back to a string
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION rdm_bit(num_bits INTEGER)
RETURNS varbit AS $$
DECLARE
    i INTEGER;
    bit_str VARCHAR := '';
BEGIN

    FOR i IN 1..num_bits LOOP
        bit_str := bit_str || (CASE WHEN random() < 0.5 THEN '0' ELSE '1' END);
    END LOOP;

    RETURN bit_str::varbit;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION rdm_varbit(num_bits INTEGER)
RETURNS varbit AS $$
DECLARE
    i INTEGER;
    bit_str VARCHAR := '';
BEGIN
    num_bits = ceil(random() * num_bits);

    FOR i IN 1..num_bits LOOP
        bit_str := bit_str || (CASE WHEN random() < 0.5 THEN '0' ELSE '1' END);
    END LOOP;

    RETURN bit_str::varbit;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION rdm_decimal(total_digits INTEGER, decimal_digits INTEGER)
RETURNS NUMERIC AS $$
DECLARE
    int_part INTEGER;
    decimal_part NUMERIC;
    max_value INTEGER;
    factor NUMERIC;
BEGIN
    total_digits = ceil(random() * total_digits);
    decimal_digits = ceil(total_digits - random() * decimal_digits);

    -- Calculate the integer part's maximum value
    max_value := 10 ^ (total_digits - decimal_digits) - 1;

    -- Generate a random integer part
    int_part := FLOOR(random() * max_value);

    -- Calculate the factor for the decimal part based on scale
    factor := 10 ^ decimal_digits;

    -- Generate a random decimal part
    decimal_part := FLOOR(random() * factor) / factor;

    -- Combine the integer and decimal parts
    RETURN int_part + decimal_part;
END;
$$ LANGUAGE plpgsql;


CREATE TABLE my_table (
    random_ascii_char char(32),
    random_unicode_char char(32),
    empty_char char(32),
    null_char char(32),
    random_ascii_varchar varchar,
    random_unicode_varchar varchar,
    empty_varchar varchar,
    null_varchar varchar,
    random_ascii_text text,
    random_unicode_text text,
    empty_text text,
    null_text text,
    random_int8 int8,
    null_int8 int8,
    random_int4 int4,
    null_int4 int4,
    random_int2 int2,
    null_int2 int2,
    random_float8 float8,
    null_float8 float8,
    random_float4 float4,
    null_float4 float4,
    random_numeric numeric(10, 5),
    null_numeric numeric(10, 5),
    random_timestamp timestamp,
    null_timestamp timestamp,
    random_timestamptz timestamptz,
    null_timestamptz timestamptz,
    random_date date,
    null_date date,
    random_time time,
    null_time time,
    random_interval interval,
    null_interval interval,
    random_bytea bytea,
    empty_bytea bytea,
    null_bytea bytea,
    random_box box,
    null_box box,
    random_circle circle,
    null_circle circle,
    random_line line,
    null_line line,
    random_lseg lseg,
    null_lseg lseg,
    random_path path,
    null_path path,
    random_point point,
    null_point point,
    random_polygon polygon,
    null_polygon polygon,
    random_bool bool,
    null_bool bool,
    random_bit bit(32),
    null_bit bit(32),
    random_varbit varbit,
    null_varbit varbit,
    random_uuid uuid,
    null_uuid uuid,
    random_ascii_json json,
    random_unicode_json json,
    null_json json,
    random_ascii_jsonb jsonb,
    random_unicode_jsonb jsonb,
    null_jsonb jsonb,
    random_inet inet,
    null_inet inet,
    random_macaddr macaddr,
    null_macaddr macaddr,
    random_cidr cidr,
    null_cidr cidr,
    random_ascii_xml xml,
    random_unicode_xml xml,
    null_xml xml,
    random_timetz timetz,
    null_timetz timetz,
    random_pg_lsn pg_lsn,
    null_pg_lsn pg_lsn,
    random_ascii_tsquery tsquery,
    random_unicode_tsquery tsquery,
    null_tsquery tsquery,
    random_ascii_tsvector tsvector,
    random_unicode_tsvector tsvector,
    null_tsvector tsvector,
    normal_serial serial,
    normal_bigserial bigserial
);

DO $$ 
DECLARE 
    counter INTEGER := 0;
BEGIN 
    WHILE counter < 1000 LOOP
        INSERT INTO my_table (
            random_ascii_char,
            random_unicode_char,
            empty_char,
            null_char,
            random_ascii_varchar,
            random_unicode_varchar,
            empty_varchar,
            null_varchar,
            random_ascii_text,
            random_unicode_text,
            empty_text,
            null_text,
            random_int8,
            null_int8,
            random_int4,
            null_int4,
            random_int2,
            null_int2,
            random_float8,
            null_float8,
            random_float4,
            null_float4,
            random_numeric,
            null_numeric,
            random_timestamp,
            null_timestamp,
            random_timestamptz,
            null_timestamptz,
            random_date,
            null_date,
            random_time,
            null_time,
            random_interval,
            null_interval,
            random_bytea,
            empty_bytea,
            null_bytea,
            random_box,
            null_box,
            random_circle,
            null_circle,
            random_line,
            null_line,
            random_lseg,
            null_lseg,
            random_path,
            null_path,
            random_point,
            null_point,
            random_polygon,
            null_polygon,
            random_bool,
            null_bool,
            random_bit,
            null_bit,
            random_varbit,
            null_varbit,
            random_uuid,
            null_uuid,
            random_ascii_json,
            random_unicode_json,
            null_json,
            random_ascii_jsonb,
            random_unicode_jsonb,
            null_jsonb,
            random_inet,
            null_inet,
            random_macaddr,
            null_macaddr,
            random_cidr,
            null_cidr,
            random_ascii_xml,
            random_unicode_xml,
            null_xml,
            random_timetz,
            null_timetz,
            random_pg_lsn,
            null_pg_lsn,
            random_ascii_tsquery,
            random_unicode_tsquery,
            null_tsquery,
            random_ascii_tsvector,
            random_unicode_tsvector,
            null_tsvector
        )
        VALUES 
        (
            rdm_ascii(32), -- random ascii char
            rdm_unicode(32), -- random unicode char
            '', -- empty char
            null, -- null char
            rdm_ascii(32), -- random ascii varchar
            rdm_unicode(32), -- random unicode varchar
            '', -- empty varchar
            null, -- null varchar
            rdm_ascii(1000), -- random ascii text
            rdm_unicode(1000), -- random unicode text
            '', -- empty text
            null, -- null text
            floor(random() * 9223372036854775807)::int8, -- random int8
            null, -- null int8
            floor(random() * 2147483647)::int4, -- random int4
            null, -- null int4
            floor(random() * 32767)::int2, -- random int2
            null, -- null int2
            random()::float8, -- random float8
            null, -- null float8
            random()::float4, -- random float4
            null, -- null float4
            rdm_decimal(10, 5), -- random numeric
            null, -- null numeric
            NOW() - '1 year'::interval * random(), -- random timestamp
            null, -- null timestamp
            NOW() - '1 year'::interval * random(), -- random timestamptz
            null, -- null timestamptz
            CURRENT_DATE - floor(random() * 365)::int, -- random date
            null, -- null date
            CURRENT_TIME - floor(random() * 86400)::int * '1 second'::interval, -- random time
            null, -- null time
            random() * '1 year'::interval * 3, -- random interval
            null, -- null interval
            decode(md5(random()::text), 'hex'), -- random bytea
            E''::bytea, -- empty bytea
            null, -- null bytea
            box(point(0, 0), point(random()*100, random()*100)), -- random_box
            null, -- null_box
            circle(point(0, 0), random()*50), -- random_circle
            null, -- null_circle
            line(point(0, 0), point(random()*100, random()*100)), -- random_line
            null, -- null_line
            lseg(point(0, 0), point(random()*100, random()*100)), -- random_lseg
            null, -- null_lseg
            path('[(0,0), (' || random()*100 || ',' || random()*100 || '), (' || random()*100 || ',' || random()*100 || ')]'), -- random_path
            null, -- null_path
            point(random()*100, random()*100), -- random_point
            null, -- null_point 
            polygon('((' || random() * 100 || ',' || random() * 100 || '),(' || random() * 100 || ',' || random() * 100 || '),(' || random() * 100 || ',' || random() * 100 || '))'), -- random_polygon
            null, -- null_polygon
            CASE WHEN random() > 0.5 THEN true ELSE false END, -- random_bool
            null, -- null_bool
            rdm_bit(32), -- random_bit
            null, -- null_bit
            rdm_varbit(32), -- random_varbit
            null, -- null_varbit
            uuid_generate_v4(), -- random_uuid
            null, -- null_uuid
            json_build_object(rdm_ascii(32), rdm_ascii(32)), -- random_ascii_json
            json_build_object(rdm_unicode(32), rdm_unicode(32)), -- random_unicode_json
            null, -- null_json
            jsonb_build_object(rdm_ascii(32), rdm_ascii(32)), -- random_ascii_jsonb
            jsonb_build_object(rdm_unicode(32), rdm_unicode(32)), -- random_unicode_jsonb
            null, -- null_jsonb
            inet(floor(random()*255)::int || '.' || floor(random()*255)::int || '.' || floor(random()*255)::int || '.' || floor(random()*255)::int), -- random_inet
            null, -- null_inet
            macaddr(lpad(to_hex(floor(random()*256)::int), 2, '0') || ':' || lpad(to_hex(floor(random()*256)::int), 2, '0') || ':' ||  lpad(to_hex(floor(random()*256)::int), 2, '0') || ':' ||  lpad(to_hex(floor(random()*256)::int), 2, '0') || ':' ||  lpad(to_hex(floor(random()*256)::int), 2, '0') || ':' ||  lpad(to_hex(floor(random()*256)::int), 2, '0')), -- random_macaddr
            null, -- null_macaddr
            cidr(floor(random()*255)::int || '.' || floor(random()*255)::int || '.' || floor(random()*255)::int || '.0/24'), -- random_cidr
            null, -- null_cidr
            xmlelement(name foo, xmlattributes(rdm_ascii(32) as bar), rdm_ascii(32)), -- random_ascii_xml
            xmlelement(name foo, xmlattributes(rdm_unicode(32) as bar), rdm_unicode(32)), -- random_unicode_xml
            null, -- null_xml
            CURRENT_TIME - floor(random()*86400)::int * '1 second'::interval, -- random_timetz
            null, -- null_timetz
            (lpad(to_hex(floor(random() * 16777215)::integer), 8, '0') || '/' || lpad(to_hex(floor(random() * 16777215)::integer), 8, '0'))::pg_lsn, -- random_pg_lsn
            null, -- null_pg_lsn
            plainto_tsquery(rdm_ascii(32)), -- random_ascii_tsquery
            plainto_tsquery(rdm_unicode(32)), -- random_unicode_tsquery
            null, -- null_tsquery
            to_tsvector(rdm_ascii(32)), -- random_ascii_tsvector
            to_tsvector(rdm_unicode(32)), -- random_unicode_tsvector
            null -- null_tsvector
        );
        counter := counter + 1;
    END LOOP;
END;
$$ LANGUAGE plpgsql;