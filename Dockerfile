FROM postgres:15-alpine
COPY ../migrations/init.sql /docker-entrypoint-initdb.d/
EXPOSE 5432