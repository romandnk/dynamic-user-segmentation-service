CREATE TABLE segments (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(255) NOT NULL UNIQUE,
    auto_add_percentage DECIMAL(5, 2) NOT NULL
);

CREATE INDEX idx_segments_slug ON segments (slug);

CREATE TABLE user_segments (
    user_id INTEGER NOT NULL,
    segment_id INTEGER NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    operation VARCHAR(6) NOT NULL,
    deleted BOOLEAN DEFAULT FALSE
);

CREATE INDEX idx_user_segments_user_id ON user_segments (user_id);
CREATE INDEX idx_user_segments_segment_id ON user_segments (segment_id);

