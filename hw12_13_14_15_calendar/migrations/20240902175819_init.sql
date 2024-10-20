-- +goose Up
-- +goose StatementBegin
CREATE table Events
(
    id          text primary key,
    title       text,
    date_time   timestamptz not null default now(),
    duration    bigint,
    description text,
    owner_id     text,
    notify_time  bigint,
    notified     boolean default false
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table Events;
-- +goose StatementEnd
