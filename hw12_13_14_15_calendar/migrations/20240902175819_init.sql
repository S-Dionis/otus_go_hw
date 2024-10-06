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
    notify_time  bigint
);

CREATE table Notifications
(
    id       text primary key,
    title    text,
    date     timestamptz not null default now(),
    "user"     text
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table Events;
DROP table Notifications;
-- +goose StatementEnd
