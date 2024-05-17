#!/bin/bash
TIME=`date +%Y%m%d%H%M%S`
mysqldump --column-statistics=0 -h ${HOST} -P${PORT} -uroot -p${PASSWORD} --databases ${DATABASE} > /tmp/${DATABASE}_${TIME}.sql
mc cp /tmp/${DATABASE}_${TIME}.sql s3/backups/mysql/${APPLICATION}/${HOST}/
