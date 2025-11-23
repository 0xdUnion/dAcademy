CREATE TABLE IF NOT EXISTS users
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    username          TEXT NOT NULL UNIQUE,
    password          TEXT NOT NULL,
    join_time         DATETIME DEFAULT CURRENT_TIMESTAMP,
    two_factor_secret TEXT
);

CREATE TABLE IF NOT EXISTS sessions
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id     TEXT NOT NULL,
    session     TEXT NOT NULL,
    create_time DATETIME DEFAULT CURRENT_TIMESTAMP
);