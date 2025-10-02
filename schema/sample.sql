USE app_db;

CREATE TABLE IF NOT EXISTS orders (
    id INT AUTO_INCREMENT PRIMARY KEY,
    customer VARCHAR(100) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO orders (customer, amount, status) VALUES
('Alice', 120.50, 'PENDING'),
('Bob', 87.30, 'PAID'),
('Charlie', 42.00, 'CANCELLED');

UPDATE orders SET status = 'SHIPPED' WHERE customer = 'Alice';
DELETE FROM orders WHERE customer = 'Charlie';

SELECT * FROM orders;
