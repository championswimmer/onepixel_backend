create database metabase;
create user metabase with encrypted password 'metabase';
grant all privileges on database metabase to metabase;

--- connect to metabase
grant all privileges on schema public to metabase;
grant all privileges on all tables in schema public to metabase;