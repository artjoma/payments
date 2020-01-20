CREATE TABLE IF NOT EXISTS accounts
(
    id          SERIAL8         PRIMARY KEY,
    balance     NUMERIC(24, 6)  DEFAULT 0   NOT NULL,
    created_at  TIMESTAMP(3)    DEFAULT (now() AT TIME ZONE 'UTC') NOT NULL,
    updated_at  TIMESTAMP(3)    DEFAULT (now() AT TIME ZONE 'UTC') NOT NULL
);