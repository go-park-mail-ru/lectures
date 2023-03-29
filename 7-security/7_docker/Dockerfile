FROM golang:1.19.0
COPY ./code /go/src/my_super_app
WORKDIR /go/src/my_super_app
RUN go install .

EXPOSE 8080/tcp

CMD [ "my_super_app" ]