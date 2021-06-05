FROM golang:1.16-alpine
COPY ./migrations /var/tern/migrations
COPY ./tern-container.conf /var/tern/tern.conf
WORKDIR /var/tern

RUN go get github.com/jackc/tern

CMD ["tern", "migrate", "--migrations", "./migrations"]