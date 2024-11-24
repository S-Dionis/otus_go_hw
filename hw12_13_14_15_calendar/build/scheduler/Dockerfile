FROM golang:1.22 as build

ENV BIN_FILE /opt/scheduler/scheduler-app
ENV CODE_DIR /go/src/scheduler

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
        -o ${BIN_FILE} cmd/scheduler/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="scheduler"
LABEL MAINTAINERS="student@otus.ru"

ENV BIN_FILE "/opt/scheduler/scheduler-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/scheduler/scheduler_config.yaml
ENV MIGRATION_FILE /etc/scheduler/20240902175819_init.sql
COPY ./configs/deploy/scheduler_config.yaml ${CONFIG_FILE}
COPY ./migrations/20240902175819_init.sql ${MIGRATION_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
