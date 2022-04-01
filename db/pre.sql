-- pre
-- drop user if exists superjer;
-- create user superjer with password 'superjer';
-- alter user superjer with superuser;
drop database if exists superjerlab;
create database superjerlab;
alter database superjerlab owner to superjer;