CREATE TABLE
    villas (
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) UNIQUE NOT NULL,
        description VARCHAR(255) NOT NULL,
        category_id INT,
        location_id INT,
        min_guest INT NOT NULL DEFAULT 1,
        bedrooms INT NOT NULL DEFAULT 1,
        price FLOAT NOT NULL,
        baths INT NOT NULL DEFAULT 1,
        image_urls JSON NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        FOREIGN KEY (category_id) REFERENCES categories (id) ON UPDATE SET NULL ON DELETE SET NULL,
        FOREIGN KEY (location_id) REFERENCES locations (id) ON UPDATE SET NULL ON DELETE SET NULL
    );