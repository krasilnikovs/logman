CREATE TABLE servers (
    `id` INTEGER PRIMARY KEY,
    `name` TEXT NOT NULL,
    `host` TEXT NOT NULL,
    `log_location_path` TEXT NOT NULL,
    `log_location_format` TEXT NOT NULL,
    `created_at` TEXT NOT NULL,
    `updated_at` TEXT NOT NULL
);