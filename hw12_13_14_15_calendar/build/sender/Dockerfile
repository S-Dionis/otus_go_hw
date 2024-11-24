FROM golang:1.22 as build

ENV BIN_FILE /opt/sender/sender-app
ENV CODE_DIR /go/src/sender

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/sender/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="sender"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/sender/sender-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/sender/sender_config.yaml
ENV MIGRATION_FILE /etc/sender/20240902175819_init.sql
COPY ./configs/deploy/sender_config.yaml ${CONFIG_FILE}
COPY ./migrations/20240902175819_init.sql ${MIGRATION_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
