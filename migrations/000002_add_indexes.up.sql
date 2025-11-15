-- Индекс для быстрого пользователей команды
CREATE INDEX idx_users_team_id ON users(team_id);

-- Индекс для фильтрации активных пользователей
CREATE INDEX idx_users_is_active ON users(is_active);

-- Cоставной индекс для поиска активных пользователей команды
CREATE INDEX idx_users_team_active ON users(team_id, is_active);

-- Индексы для фильтрации PR по статусу и автору
CREATE INDEX idx_pr_status_id ON pull_requests(status_id);
CREATE INDEX idx_pr_author_id ON pull_requests(author_id);

-- Индекс для быстрого поиска PR пользователя
CREATE INDEX idx_pr_reviewers_user_id ON pr_reviewers(user_id);

-- Индекс для сортировки по времени создания (для статистики)
CREATE INDEX idx_pr_created_at ON pull_requests(created_at DESC);
