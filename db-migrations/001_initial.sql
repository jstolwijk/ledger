create table journal_entry(
    id uuid primary key not null,
    idempotency_key text not null,
    from_account text not null,
    to_account text not null,
    amount_value bigint not null,
    amount_currency text not null,
    metadata jsonb not null
);
---- create above / drop below ----