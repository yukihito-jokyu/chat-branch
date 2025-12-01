-- +goose Up
CREATE TABLE chat_summary (
    offset bigint NOT NULL AUTO_INCREMENT,
    uuid varchar(36) NOT NULL,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    payload blob DEFAULT NULL,
    metadata json DEFAULT NULL,
    topic varchar(255) NOT NULL,
    PRIMARY KEY (offset)
);

-- +goose Down
DROP TABLE chat_summary;
