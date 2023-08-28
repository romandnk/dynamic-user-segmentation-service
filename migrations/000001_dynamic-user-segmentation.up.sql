CREATE TABLE segments (
    slug VARCHAR(255) PRIMARY KEY,
    auto_add_percentage SMALLINT NOT NULL
);

CREATE TABLE user_segments (
    user_id INTEGER NOT NULL,
    segment_slug VARCHAR(255),
    FOREIGN KEY (segment_slug) REFERENCES segments (slug) ON DELETE NO ACTION
);

CREATE INDEX idx_user_segments_user_id ON user_segments (user_id);

CREATE TABLE operations (
    user_id INTEGER NOT NULL,
    segment_slug VARCHAR(255) NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    operation VARCHAR(6) NOT NULL,
    auto_add BOOLEAN NOT NULL
);

CREATE INDEX idx_operations_user_id ON operations (user_id);
