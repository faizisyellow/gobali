ALTER TABLE user_invitation DROP FOREIGN KEY user_fk;
ALTER TABLE user_invitation ADD FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;