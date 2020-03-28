#!/bin/sh
docker-compose exec postgres bash -c "chmod 0775 docker-entrypoint-initdb.d/init-database.sh"
docker-compose exec postgres bash -c "./docker-entrypoint-initdb.d/init-database.sh"