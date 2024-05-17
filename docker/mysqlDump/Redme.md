# mysql backup database image

```dockerfile
# ubuntu size 28.17 MB
FROM ubuntu:22.04

WORKDIR /

# copy files
COPY ./entrypoint.sh /
COPY ./mc /usr/local/sbin/

# execute instructions
# mc configuration customize
RUN chmod +x /usr/local/sbin/mc && \
	chmod +x /entrypoint.sh && \
	apt update && \
	apt install mysql-client -y && \
	mc config host add s3 http://192.168.1.95:9000 PGd0IEaLCrVGNg1ZLozC CSPLPfGwcVeVu7Vrl9lKDu5jx4gnXau0fy5YC7Z9 --api s3v4

# debug run docker run -it [images:target] /bin/bash
CMD ["/entrypoint.sh"]
```

## scripts

### environment variable
```bash
HOST: database address
PASSWORD: database user(root) password
DATABASE: database name
PORT: database port
APPLICATION: customize
```

## minio
1. create Access key
2. bucket name backups
3. directory: mysql/[APPLICATION]/[HOST]/[DATABASE]_[TIME].sql

## Tips
```shell
ERROR: The difference between the request time and the server‘s time is too large.
Resolve: -v /etc/localtime:/etc/localtime:ro
```