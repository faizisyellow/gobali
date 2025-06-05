CREATE TABLE users(
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(16) NOT NULL,
    email VARCHAR(28) NOT NULL UNIQUE,
    password BINARY(40) NOT NULL,
    is_active BOOLEAN DEFAULT 0,
    role_id INT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES  roles(id)    
);