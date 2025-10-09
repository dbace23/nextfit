CREATE DATABASE IF NOT EXISTS nextfit;
USE nextfit;


-- =========================
-- Users
-- =========================
CREATE TABLE IF NOT EXISTS users (
    user_id         INT AUTO_INCREMENT PRIMARY KEY,
    full_name       VARCHAR(255) NOT NULL,
    password_hash   VARCHAR(255) NOT NULL,
    email           VARCHAR(255),
    isAdmin        TINYINT NOT NULL DEFAULT 1,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      TIMESTAMP NULL DEFAULT NULL,
    UNIQUE KEY uq_users_email (email)
);

-- =========================
-- Categories
-- =========================
CREATE TABLE IF NOT EXISTS categories (
    category_id     INT AUTO_INCREMENT PRIMARY KEY,
    category_name   VARCHAR(255) NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NULL DEFAULT NULL ,
    deleted_at      TIMESTAMP NULL DEFAULT NULL,
    UNIQUE KEY uq_categories_name (category_name)
);

-- =========================
-- Products
-- =========================
CREATE TABLE IF NOT EXISTS products (
    product_id      INT AUTO_INCREMENT PRIMARY KEY,
    product_name    VARCHAR(255) NOT NULL,
    category_id     INT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    selling_price   DECIMAL(12,2) NOT NULL,
    updated_at      TIMESTAMP NULL DEFAULT NULL,
    deleted_at      TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_products_category
        FOREIGN KEY (category_id)
        REFERENCES categories (category_id)
        ON DELETE RESTRICT
);


-- ⚠️ NOTE: MySQL does not support partial unique indexes.
-- You must enforce "one default address per user" in application logic or with a trigger.

-- =========================
-- Orders
-- =========================
CREATE TABLE IF NOT EXISTS orders (
    order_id        INT AUTO_INCREMENT PRIMARY KEY,
    user_id         INT,
    address         VARCHAR(100),
    paid_amount     DECIMAL(12,2) NOT NULL,
    shipping_fee    DECIMAL(12,2) NOT NULL,
    payment_type    ENUM('CASH', 'TRANSFER', 'EWALLET', 'CREDIT_CARD') NOT NULL,
    order_status    ENUM('PENDING', 'COMPLETED') NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NULL DEFAULT NULL,
    deleted_at      TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_orders_user
        FOREIGN KEY (user_id)
        REFERENCES users (user_id)
        ON DELETE RESTRICT
);

-- =========================
-- Order Details
-- =========================
CREATE TABLE IF NOT EXISTS order_details (
    order_detail_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id        INT NOT NULL,
    product_id      INT NOT NULL,
    quantity        INT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP NULL DEFAULT NULL,
    deleted_at      TIMESTAMP NULL DEFAULT NULL,
    CONSTRAINT fk_order_details_order
        FOREIGN KEY (order_id)
        REFERENCES orders (order_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_order_details_product
        FOREIGN KEY (product_id)
        REFERENCES products (product_id)
        ON DELETE RESTRICT
);


INSERT INTO categories (category_name, created_at, updated_at, deleted_at) VALUES
('Soccer',       '2025-10-01 08:00:00', '2025-10-05 10:00:00', NULL),
('Basketball',   '2025-10-02 09:00:00', NULL,                 NULL),
('Baseball',     '2025-09-28 07:30:00', '2025-10-03 15:00:00', NULL),
('Tennis',       '2025-10-01 12:00:00', '2025-10-07 08:45:00', '2025-10-08 14:00:00'),
('Swimming',     '2025-10-03 11:00:00', NULL,                 NULL),
('Volleyball',   '2025-10-04 13:00:00', '2025-10-06 09:30:00', NULL),
('Table Tennis', '2025-10-05 14:00:00', '2025-10-08 10:15:00', '2025-10-09 11:00:00'),
('Badminton',    '2025-09-30 10:30:00', NULL,                 NULL),
('Golf',         '2025-10-06 15:00:00', '2025-10-09 16:45:00', NULL),
('Rugby',        '2025-09-29 08:15:00', NULL,                 '2025-10-07 13:30:00');
