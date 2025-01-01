CREATE TABLE google_users
(
    id             UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    google_id      VARCHAR(256) NOT NULL UNIQUE,
    email          TEXT         NOT NULL UNIQUE,
    verified_email BOOLEAN      NOT NULL DEFAULT FALSE,
    name           TEXT,
    picture        TEXT,
    locale         TEXT
);