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
