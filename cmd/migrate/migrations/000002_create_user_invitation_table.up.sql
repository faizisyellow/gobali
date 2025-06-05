CREATE TABLE
    user_invitation (
        token BINARY(32) NOT NULL,
        user_id INT NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );