FROM golang

RUN mkdir -p /go/src/tournament
WORKDIR /go/src/tournament

COPY . /go/src/tournament
RUN go get -d -v
RUN go build -o main .

EXPOSE 9000
CMD [ "/go/src/tournament/main", "-env", "local" ]
