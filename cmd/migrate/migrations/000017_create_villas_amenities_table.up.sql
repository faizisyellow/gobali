CREATE TABLE IF NOT EXISTS
    villas_amenities (
        villa_id INT NOT NULL,
        amentity_id INT NOT NULL,
        PRIMARY KEY (villa_id, amentity_id),
        FOREIGN KEY (villa_id) REFERENCES villas (id) ON DELETE CASCADE,
        FOREIGN KEY (amentity_id) REFERENCES amenities (id) ON DELETE CASCADE
    )