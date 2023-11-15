CREATE TABLE log_infos (
    `id` INT PRIMARY KEY,
    `name` TEXT NOT NULL,
    `location` TEXT NOT NULL,
    `format` TEXT NOT NULL,
    `created_at` TEXT NOT NULL,
    `updated_at` TEXT NOT NULL
);

CREATE TABLE servers (
    `id` INT PRIMARY KEY,
    'log_infos_id' INT NOT NULL,
    `name` TEXT NOT NULL,
    `host` TEXT NOT NULL,
    `created_at` TEXT NOT NULL,
    `updated_at` TEXT NOT NULL,
    FOREIGN KEY(log_infos_id) REFERENCES log_infos(id)
);