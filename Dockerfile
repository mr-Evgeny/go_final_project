FROM golang:1.22 as go

WORKDIR /app

# подготавливаем зависимости
COPY go.mod go.sum ./ 
RUN go mod download

#копируем программу и БД
COPY . .

ENV TODO_DBFILE /app/db/scheduler.db

RUN GOOS=linux GOARCH=amd64 go test
RUN GOOS=linux GOARCH=amd64 go build -o todo

# Запускаем
ENTRYPOINT ["/app/todo"]
