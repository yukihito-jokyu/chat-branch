-- +goose Up
-- 1. カラムの追加 (NULL許容)
ALTER TABLE messages 
ADD COLUMN parent_message_uuid CHAR(255) NULL AFTER uuid;

-- 2. 外部キー制約の追加 (自己参照)
-- ON DELETE CASCADE: 親メッセージが消えたらリプライも削除する
ALTER TABLE messages
ADD CONSTRAINT fk_messages_parent
FOREIGN KEY (parent_message_uuid) 
REFERENCES messages(uuid)
ON DELETE CASCADE
ON UPDATE CASCADE;

-- +goose Down
-- 1. 外部キー制約の削除
ALTER TABLE messages 
DROP FOREIGN KEY fk_messages_parent;

-- 2. カラムの削除
ALTER TABLE messages 
DROP COLUMN parent_message_uuid;
