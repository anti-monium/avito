#!/bin/bash

docker-compose down --volumes
docker-compose up --build -d
sleep 20

echo "Создание таблиц в базе данных"
docker exec -i $(docker ps -qf "name=postgres") psql -U postgres -d postgres < pkg/database/sql/apartment.down.sql
docker exec -i $(docker ps -qf "name=postgres") psql -U postgres -d postgres < pkg/database/sql/apartment.up.sql