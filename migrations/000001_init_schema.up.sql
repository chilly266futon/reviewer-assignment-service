-- Статусы pull request'ов
CREATE TABLE pr_statuses (
    id SMALLINT PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE
);

INSERT INTO pr_statuses (id, name) VALUES
    (1, 'OPEN'),
    (2, 'MERGED');


-- Команды
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Пользователи
CREATE TABLE users (
    id VARCHAR(100) PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Pull Requests
CREATE TABLE pull_requests (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    author_id VARCHAR(100) NOT NULL REFERENCES users(id),
    status_id SMALLINT NOT NULL REFERENCES pr_statuses(id) DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    merged_at TIMESTAMP
);

-- Назначения ревьюверов на pull request'ы
CREATE TABLE pr_reviewers (
    pull_request_id VARCHAR(100) NOT NULL REFERENCES pull_requests (id) ON DELETE CASCADE,
    user_id         VARCHAR(100) NOT NULL REFERENCES users (id),
    assigned_at     TIMESTAMP    NOT NULL DEFAULT NOW(),
    PRIMARY KEY (pull_request_id, user_id)
);

