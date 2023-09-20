#docker build -t webis/bigbro:20.Dec.2018 .
FROM golang:1.19

RUN go install github.com/hscells/bigbro/cmd/bigbro@latest

CMD [ "bigbro", "--filename", "my_application.logr", "csv" ]
