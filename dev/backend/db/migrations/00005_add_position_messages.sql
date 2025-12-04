-- +goose Up
-- 1. position_x カラムの追加 (context_summary の後ろに追加)
ALTER TABLE messages
ADD COLUMN position_x FLOAT NOT NULL DEFAULT 0 AFTER context_summary;

-- 2. position_y カラムの追加 (position_x の後ろに追加)
ALTER TABLE messages
ADD COLUMN position_y FLOAT NOT NULL DEFAULT 0 AFTER position_x;

-- +goose Down
-- 1. カラムの削除
ALTER TABLE messages
DROP COLUMN position_y;

ALTER TABLE messages
DROP COLUMN position_x;