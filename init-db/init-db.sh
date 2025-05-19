#!/bin/bash
set -e

SQL_COMMANDS=$(cat <<-END
USE ${MYSQL_DATABASE};

CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO tasks (title, description, status) VALUES
('Setup Docker environment', 'Configure Dockerfile and Docker Compose for the project', 'completed'),
('Implement gRPC User Service', 'Create gRPC service for user management', 'pending'),
('Write unit tests', 'Add unit tests for critical components', 'todo')
ON DUPLICATE KEY UPDATE title=VALUES(title);

END
)

echo "Running database initialization script..."
mysql -u root -p"${MYSQL_ROOT_PASSWORD}" "${MYSQL_DATABASE}" <<-EOSQL
    ${SQL_COMMANDS}
EOSQL

echo "Database initialization script finished."