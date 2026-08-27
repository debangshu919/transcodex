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
