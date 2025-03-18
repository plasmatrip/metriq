#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$postgres" --dbname "$password" <<-EOSQL
    CREATE USER metriq WITH ENCRYPTED PASSWORD 'password';
    CREATE DATABASE metriq OWNER 'metriq';
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO metriq;
EOSQL