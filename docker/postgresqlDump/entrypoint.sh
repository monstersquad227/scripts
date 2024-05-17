#!/bin/bash
TIME=`date +%Y%m%d%H%M%S`
pg_dump -h ${HOST} -U ${USERNAME} ${DATABASE} > /tmp/${DATABASE}_${TIME}.sql
mc cp /tmp/${DATABASE}_${TIME}.sql s3/backups/postgresql/${APPLICATION}/${HOST}/