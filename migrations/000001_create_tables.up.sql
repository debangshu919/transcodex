CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    name TEXT NOT NULL,

    prefix TEXT NOT NULL,
    hash TEXT NOT NULL UNIQUE,

    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    storage_key TEXT NOT NULL UNIQUE,

    filename TEXT NOT NULL,
    content_type TEXT,
    extension TEXT,

    size_bytes BIGINT NOT NULL,

    status TEXT NOT NULL,

    metadata JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS conversions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    input_file_id UUID NOT NULL
        REFERENCES files(id),

    output_file_id UUID
        REFERENCES files(id),

    source_format TEXT NOT NULL,
    target_format TEXT NOT NULL,

    status TEXT NOT NULL,

    options JSONB,

    progress REAL NOT NULL DEFAULT 0,

    error_code TEXT,
    error_message TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    conversion_id UUID NOT NULL
        REFERENCES conversions(id)
        ON DELETE CASCADE,

    status TEXT NOT NULL,

    attempt INTEGER NOT NULL DEFAULT 1,

    worker_id TEXT,

    error_code TEXT,
    error_message TEXT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS credit_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL UNIQUE
        REFERENCES users(id)
        ON DELETE CASCADE,

    balance BIGINT NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (balance >= 0)
);

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
