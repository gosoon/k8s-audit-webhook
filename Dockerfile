FROM daocloud.io/library/golang
MAINTAINER tianfeiyu 
RUN mkdir -p /home/audit-webhook
COPY ./main /home/audit-webhook
EXPOSE 8081
CMD ["/home/audit-webhook/main"]
