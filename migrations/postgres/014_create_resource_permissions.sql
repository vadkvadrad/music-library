-- Создание таблицы resource_permissions
CREATE TABLE IF NOT EXISTS resource_permissions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    permission VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_resource_permissions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_resource_permissions_user_id ON resource_permissions(user_id);
CREATE INDEX IF NOT EXISTS idx_resource_permissions_resource_id ON resource_permissions(resource_id);

