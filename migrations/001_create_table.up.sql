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
    author_id INTEGER REFERENCES users(id) DEFAULT 0,
    assigned_id INTEGER REFERENCES users(id) DEFAULT 0,
    title TEXT NOT NULL,
    content TEXT
);

CREATE TABLE tasks_labels (
    task_id INTEGER REFERENCES tasks(id) ON DELETE CASCADE,
    label_id INTEGER REFERENCES labels(id) ON DELETE CASCADE
);

INSERT INTO users (id, name)
VALUES (0, 'John Doe'), (1, 'Jane Doe');

INSERT INTO labels (id, name)
VALUES (0, 'Work'), (1, 'Personal');

INSERT INTO tasks (id, author_id, title, content)
VALUES (0, 0, 'Example task #1', 'This is a test task.'),
       (1, 0, 'Example task #2', 'This is another test task.');

INSERT INTO tasks_labels (task_id, label_id)
VALUES (0, 1), (1, 1)
