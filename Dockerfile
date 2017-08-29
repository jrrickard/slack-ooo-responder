FROM alpine:latest
RUN apk add -U tzdata
RUN cp /usr/share/zoneinfo/America/Denver /etc/localtime
RUN echo "America/Denver" >  /etc/timezone
RUN apk --no-cache add ca-certificates && update-ca-certificates
COPY bin/linux_amd64/slack-ooo-responder /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/slack-ooo-responder"]
