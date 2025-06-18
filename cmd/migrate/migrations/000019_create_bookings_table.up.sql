CREATE TABLE bookings(
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    villa_id INT,
    start_at DATE NOT NULL,
    end_at DATE NOT NULL,
    status ENUM('open','complete','expire','cancel') NOT NULL DEFAULT 'open',
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255),
    expire_at TIME,
    villa_name VARCHAR(255) NOT NULL,
    villa_location VARCHAR(255) NOT NULL,
    villa_price INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(villa_id) REFERENCES villas(id) ON DELETE SET NULL
)