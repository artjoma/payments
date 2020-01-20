CREATE TABLE IF NOT EXISTS transactions
(
    id          SERIAL8         PRIMARY KEY,
    tx_id       TEXT            NOT NULL,
    created_at  TIMESTAMP(3)    DEFAULT (now() AT TIME ZONE 'UTC') NOT NULL,
    prev_balance NUMERIC(24, 6)  DEFAULT 0 NOT NULL,
    amount      NUMERIC(24, 6)  DEFAULT 0 NOT NULL,
    balance     NUMERIC(24, 6)  DEFAULT 0 NOT NULL,
    state       TEXT            NOT NULL,
    source      TEXT            NOT NULL,
    status      TEXT            NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS TX_ID_INDEX ON transactions (tx_id);