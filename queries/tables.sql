
CREATE TABLE IF NOT EXISTS episodes (
    id TEXT PRIMARY KEY,
    value TEXT,
    transcription_id INTEGER,
    summary_id INTEGER,
    source_type_id INTEGER,
    
    FOREIGN KEY (transcription_id) REFERENCES transcriptions(id)
    FOREIGN KEY (summary_id) REFERENCES summaries(id),
    FOREIGN KEY (source_type_id) REFERENCES source_types(id)
);

CREATE TABLE IF NOT EXISTS transcriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    value BLOB
);

CREATE TABLE IF NOT EXISTS summaries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    value BLOB
);

CREATE TABLE IF NOT EXISTS source_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL
);

-- Insert default values for source_types
INSERT OR IGNORE INTO source_types (id, name) VALUES (1, 'youtube');
INSERT OR IGNORE INTO source_types (id, name) VALUES (2, 'podcast');