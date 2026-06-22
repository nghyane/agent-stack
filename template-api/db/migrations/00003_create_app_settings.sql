-- +goose Up
-- +goose StatementBegin
-- Singleton row holding runtime-mutable operational settings as JSONB.
-- This is NOT for secrets or wiring (those stay in env). It is for operational
-- toggles that change without a redeploy (e.g. which AI model is active).
CREATE TABLE app_settings (
    id         INT PRIMARY KEY DEFAULT 1,
    data       JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT app_settings_singleton CHECK (id = 1)
);

INSERT INTO app_settings (id) VALUES (1) ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE app_settings;
-- +goose StatementEnd
