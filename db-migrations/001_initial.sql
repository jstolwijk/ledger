create table journal_entry(
    id uuid primary key not null,
    idempotency_key text not null,
    from_account text not null,
    to_account text not null,
    amount_value bigint not null,
    amount_currency text not null,
    metadata jsonb not null
);
CREATE TABLE ledger_entry(
    id uuid primary key not null,
    ledger_id uuid not null,
    ledger_name text not null,
    ledger_version bigint not null,
    journal_entry_id uuid not null,
    currency text not null,
    from_amount bigint,
    to_amount bigint,
    ledger_sum_to bigint not null,
    ledger_sum_from bigint not null
);
CREATE UNIQUE INDEX ledger_entry_version_idx ON ledger_entry(ledger_name, ledger_version);
CREATE TABLE ledger(
    id uuid primary key not null,
    name text not null,
    version bigint not null,
    sum_from bigint not null,
    sum_to bigint not null
);
CREATE UNIQUE INDEX ledger_name_idx ON ledger(name);
---- create above / drop below ----