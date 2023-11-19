CREATE TABLE servers (
    `id` INTEGER PRIMARY KEY,
    credential_id INTEGER NOT NULL,
    `name` TEXT NOT NULL,
    `host` TEXT NOT NULL,
    `log_location_path` TEXT NOT NULL,
    `log_location_format` TEXT NOT NULL,
    `created_at` TEXT NOT NULL,
    `updated_at` TEXT NOT NULL,
    FOREIGN KEY (credential_id) REFERENCES credentials(id)
);

CREATE TABLE credentials (
    `id` INTEGER PRIMARY KEY,
    `name` TEXT NOT NULL,
    `path` TEXT NOT NULL,
    `created_at` TEXT NOT NULL,
    `updated_at` TEXT NOT NULL
);