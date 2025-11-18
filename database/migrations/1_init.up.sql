CREATE TABLE IF NOT EXISTS courses
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    slug          TEXT NOT NULL UNIQUE,
    name          TEXT NOT NULL,
    description   TEXT,
    tags          TEXT, -- JSON
    folder        TEXT,
    chapter_count INTEGER DEFAULT 0
);