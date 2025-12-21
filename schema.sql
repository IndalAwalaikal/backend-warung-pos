-- schema.sql: create database, tables, and sample queries for Warung POS

-- 1) Create database (run as root or existing user)
CREATE DATABASE IF NOT EXISTS warung_pos
CHARACTER SET utf8mb4
COLLATE utf8mb4_general_ci;

USE warung_pos;

-- 2) Users
CREATE TABLE IF NOT EXISTS users (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(150) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(30) NOT NULL DEFAULT 'kasir',
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3) Categories
CREATE TABLE IF NOT EXISTS categories (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(120) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 4) Menus
CREATE TABLE IF NOT EXISTS menus (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(200) NOT NULL,
  description TEXT,
  price DECIMAL(12,2) NOT NULL DEFAULT 0,
  category_id BIGINT UNSIGNED NULL,
  image_url VARCHAR(512),
  is_available TINYINT(1) NOT NULL DEFAULT 1,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_menus_category (category_id),
  CONSTRAINT fk_menus_category
    FOREIGN KEY (category_id) REFERENCES categories(id)
    ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 5) Transactions
CREATE TABLE IF NOT EXISTS transactions (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  total DECIMAL(14,2) NOT NULL DEFAULT 0,
  subtotal DECIMAL(14,2) NOT NULL DEFAULT 0,
  tax DECIMAL(14,2) NOT NULL DEFAULT 0,
  discount DECIMAL(14,2) NOT NULL DEFAULT 0,
  payment_method ENUM('tunai','qris') NOT NULL DEFAULT 'tunai',
  amount_paid DECIMAL(14,2) NOT NULL DEFAULT 0,
  cashier_id BIGINT UNSIGNED NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  INDEX idx_transactions_cashier (cashier_id),
  CONSTRAINT fk_transactions_cashier
    FOREIGN KEY (cashier_id) REFERENCES users(id)
    ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 6) Transaction items
CREATE TABLE IF NOT EXISTS transaction_items (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  transaction_id BIGINT UNSIGNED NOT NULL,
  menu_id BIGINT UNSIGNED NULL,
  quantity INT NOT NULL DEFAULT 1,
  price DECIMAL(12,2) NOT NULL DEFAULT 0,
  PRIMARY KEY (id),
  INDEX idx_titems_tx (transaction_id),
  INDEX idx_titems_menu (menu_id),
  CONSTRAINT fk_titems_tx
    FOREIGN KEY (transaction_id) REFERENCES transactions(id)
    ON DELETE CASCADE,
  CONSTRAINT fk_titems_menu
    FOREIGN KEY (menu_id) REFERENCES menus(id)
    ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 7) Sample inserts (AMAN JIKA DIJALANKAN BERULANG)
INSERT INTO users (name, email, password, role)
VALUES ('Admin Warung', 'admin@warung.com', '$2a$10$Z1q7...', 'admin')
ON DUPLICATE KEY UPDATE
name = VALUES(name),
password = VALUES(password),
role = VALUES(role);

INSERT INTO categories (name)
VALUES ('Makanan'), ('Minuman')
ON DUPLICATE KEY UPDATE
name = VALUES(name);

INSERT INTO menus (name, description, price, category_id, image_url, is_available)
VALUES
('Nasi Goreng', 'Nasi goreng spesial', 18000.00, 1, '', 1),
('Es Teh', 'Es teh manis', 5000.00, 2, '', 1)
ON DUPLICATE KEY UPDATE
price = VALUES(price),
is_available = VALUES(is_available);

-- 8) Useful queries
-- Get menus with category name
SELECT m.id, m.name, m.description, m.price, m.image_url, m.is_available,
       c.name AS category
FROM menus m
LEFT JOIN categories c ON m.category_id = c.id
WHERE m.is_available = 1
ORDER BY m.name;

-- Daily report example: revenue per day (today)
SELECT DATE(created_at) AS date,
       SUM(total) AS revenue,
       COUNT(*) AS transactions
FROM transactions
WHERE DATE(created_at) = CURDATE()
GROUP BY DATE(created_at);

-- 9) Advanced: top selling menu items (today)
SELECT m.id, m.name,
       SUM(ti.quantity) AS total_qty,
       SUM(ti.quantity * ti.price) AS revenue
FROM transaction_items ti
JOIN menus m ON m.id = ti.menu_id
JOIN transactions t ON t.id = ti.transaction_id
WHERE DATE(t.created_at) = CURDATE()
GROUP BY m.id, m.name
ORDER BY total_qty DESC
LIMIT 10;

-- 10) Cleanup examples (CONTOH STATIS, BUKAN PREPARED STATEMENT)
-- delete a menu with id 1
DELETE FROM menus WHERE id = 1;

-- end of schema.sql
