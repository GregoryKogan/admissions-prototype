#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    DO
    \$do\$
    BEGIN
        IF NOT EXISTS (
            SELECT
            FROM   pg_catalog.pg_roles
            WHERE  rolname = '$POSTGRES_USER') THEN

            CREATE ROLE $POSTGRES_USER WITH LOGIN SUPERUSER PASSWORD '$POSTGRES_PASSWORD';
        END IF;

        IF NOT EXISTS (
            SELECT
            FROM   pg_catalog.pg_database
            WHERE  datname = '$POSTGRES_DB') THEN

            CREATE DATABASE $POSTGRES_DB WITH OWNER = $POSTGRES_USER
            ENCODING = 'UTF8'
            LC_COLLATE = 'en_US.utf8'
            LC_CTYPE = 'en_US.utf8'
            CONNECTION LIMIT = -1;
        END IF;
    END
    \$do\$;
EOSQL