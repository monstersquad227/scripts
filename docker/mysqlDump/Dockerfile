FROM ubuntu:22.04
WORKDIR /

COPY ./entrypoint.sh /
COPY ./mc /usr/local/sbin/

USER root

RUN chmod +x /usr/local/sbin/mc && \
	chmod +x /entrypoint.sh && \
	apt update && \
	apt install mysql-client -y && \
	mc config host add s3 http://192.168.1.95:9000 PGd0IEaLCrVGNg1ZLozC CSPLPfGwcVeVu7Vrl9lKDu5jx4gnXau0fy5YC7Z9 --api s3v4

CMD ["/entrypoint.sh"]
