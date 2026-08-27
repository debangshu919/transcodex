CREATE TABLE IF NOT EXISTS credit_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    credit_account_id UUID NOT NULL
        REFERENCES credit_accounts(id),

    conversion_id UUID
        REFERENCES conversions(id),

    type TEXT NOT NULL,

    amount BIGINT NOT NULL,

    balance_after BIGINT NOT NULL,

    description TEXT,

    idempotency_key TEXT UNIQUE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
