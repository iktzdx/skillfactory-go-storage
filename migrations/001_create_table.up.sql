DROP TABLE IF EXISTS tasks_labels, tasks, labels, users;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    opened TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    closed TIMESTAMPTZ DEFAULT NULL,
    author_id INTEGER DEFAULT 1 REFERENCES users(id) ON DELETE SET DEFAULT,
    assigned_id INTEGER DEFAULT 1 REFERENCES users(id) ON DELETE SET DEFAULT,
    title TEXT NOT NULL,
    content TEXT
);

CREATE TABLE tasks_labels (
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    label_id INTEGER REFERENCES labels(id) ON DELETE CASCADE
);

INSERT INTO users (name) VALUES ('user #1'), ('user #2');

INSERT INTO labels (name) VALUES ('example');

INSERT INTO tasks (title, content)
VALUES ('example task #1', 'This is a test task.'),
        ('example task #2', 'This is another test task.');

INSERT INTO tasks_labels VALUES (1, 1), (2, 1);
