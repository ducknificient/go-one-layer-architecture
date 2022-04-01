-- connect database first and dump sql
DROP EXTENSION IF EXISTS "uuid-ossp";
CREATE EXTENSION IF not EXISTS "uuid-ossp";

ALTER DATABASE superjerlab SET timezone TO 'Asia/Jakarta';

-- create schema type
drop schema if exists types cascade;
create schema types;
alter schema types owner to superjer;
set search_path = types,public;

create type flag as enum (  
    'valid',
    'invalid',
    'active',
    'passive',
    'enable',
    'disable',
    'up',
    'down',
    'insert',
    'update',
    'delete',
    'inactive'
);

drop schema if exists smaster cascade;
create schema smaster;
alter schema smaster owner to superjer;
set search_path = smaster,public;

create table smaster.shoes(
    id uuid primary key not null default uuid_generate_v1mc(),
    seq int default null,
    name text default null,
    description text default null,
    created timestamp without time zone default now() not null,
    createdby text not null,
    createdip text default null,
    updated timestamp without time zone default now() not null,
    updatedby text not null,
    updatedip text default null,
    flag types.flag default 'insert'::types.flag
);