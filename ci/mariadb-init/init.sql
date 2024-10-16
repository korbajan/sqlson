-- For MariaDB
CREATE TABLE IF NOT EXISTS devs (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    age INT,
    email VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description TEXT NOT NULL,
    status ENUM('active', 'inactive', 'pending') DEFAULT 'active',
    priority INT CHECK (priority >= 1 AND priority <= 5),
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    dev_id INT,  -- Foreign key reference to devs
    FOREIGN KEY (dev_id) REFERENCES devs(id) ON DELETE CASCADE
);

-- Populate devs with sample data
INSERT INTO devs (name, age, email) VALUES
('Alice Smith', 30, 'alice@example.com'),
('Bob Johnson', 25, 'bob@example.com'),
('Charlie Brown', 35, 'charlie@example.com');

-- Populate tasks with sample data
INSERT INTO tasks (description, status, priority, dev_id) VALUES
('Task 1: Complete project report', 'active', 3, 1),  -- Related to Alice
('Task 2: Review pull requests', 'pending', 2, 2),   -- Related to Bob
('Task 3: Prepare for presentation', 'inactive', 1, 3); -- Related to Charlie
