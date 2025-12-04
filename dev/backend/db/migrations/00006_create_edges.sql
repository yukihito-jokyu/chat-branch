-- +goose Up
CREATE TABLE edges (
    uuid VARCHAR(255) NOT NULL COMMENT 'UUID',
    chat_uuid VARCHAR(255) NOT NULL COMMENT 'chatsテーブルのUUID',
    source_message_uuid VARCHAR(255) NOT NULL COMMENT 'エッジの元のUUID',
    target_message_uuid VARCHAR(255) NOT NULL COMMENT 'エッジの先のUUID',
    CONSTRAINT fk_edges_chat FOREIGN KEY (chat_uuid) REFERENCES chats(uuid) ON DELETE CASCADE,
    CONSTRAINT fk_edges_source_message FOREIGN KEY (source_message_uuid) REFERENCES messages(uuid) ON DELETE CASCADE,
    CONSTRAINT fk_edges_target_message FOREIGN KEY (target_message_uuid) REFERENCES messages(uuid) ON DELETE CASCADE
) COMMENT='エッジ管理テーブル';

-- +goose Down
DROP TABLE edges;
