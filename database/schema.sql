CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_name VARCHAR(1024) NOT NULL,
    s3_url VARCHAR(1024),
    file_size BIGINT NOT NULL,
    file_type VARCHAR(50),
    upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_public BOOLEAN DEFAULT FALSE
);

CREATE TABLE shared_files (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_id INT NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    s3_url VARCHAR(1024) NOT NULL,
    shared_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add an index on user_id and file_name for efficient searching by file name
CREATE INDEX idx_files_user_id_file_name
ON files (user_id, file_name);

-- Add an index on user_id and upload_date for efficient searching by upload date
CREATE INDEX idx_files_user_id_upload_date
ON files (user_id, upload_date);

-- Add an index on user_id and file_type for efficient searching by file type
CREATE INDEX idx_files_user_id_file_type
ON files (user_id, file_type);
