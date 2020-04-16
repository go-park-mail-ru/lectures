FROM golang:1.13.0

COPY ./code /go/src/docker0
RUN go install docker0

EXPOSE 8080/tcp

CMD [ "docker0" ]