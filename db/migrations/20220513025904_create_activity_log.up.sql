CREATE TABLE activity_log (
    id INT NOT NULL PRIMARY KEY,
    user_id INT NOT NULL, 
    module VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    CONSTRAINT `fk_log_user` FOREIGN KEY(user_id) REFERENCES user(id)
);