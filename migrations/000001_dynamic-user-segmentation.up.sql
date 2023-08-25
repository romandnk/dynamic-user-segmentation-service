CREATE TABLE segments (
    slug VARCHAR(255) PRIMARY KEY ,
    auto_add_percentage SMALLINT NOT NULL,
    deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_segments_slug ON segments (slug);

CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    segment_slug VARCHAR(255) NOT NULL REFERENCES segments (slug),
    date TIMESTAMPTZ NOT NULL,
    operation VARCHAR(6) NOT NULL,
    deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_user_segments_segment_slug ON users (segment_slug);
