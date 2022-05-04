CREATE TABLE IF NOT EXISTS user_role (
    id INT NOT NULL PRIMARY KEY, 
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS user_info (
    id INT NOT NULL PRIMARY KEY,
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone VARCHAR(255),
    address VARCHAR(255),
    role_id INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    CONSTRAINT `fk_user_role` FOREIGN KEY(role_id) REFERENCES user_role(id)
);

CREATE TABLE IF NOT EXISTS user (
    id INT NOT NULL PRIMARY KEY,
    user_id INT UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    CONSTRAINT `fk_auth_user` FOREIGN KEY(user_id) REFERENCES user_info(id)
);

CREATE TABLE IF NOT EXISTS product_category (
    id INT NOT NULL PRIMARY KEY, 
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS product (
    id INT NOT NULL PRIMARY KEY,
    category_id INT NOT NULL,
    name VARCHAR(255) NOT NULL, 
    price INT NOT NULL,
    qty INT NOT NULL,
    url TEXT,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    CONSTRAINT `fk_product_cat` FOREIGN KEY(category_id) REFERENCES product_category(id)
);

CREATE TABLE IF NOT EXISTS coupon (
    id INT NOT NULL PRIMARY KEY, 
    code VARCHAR(255) NOT NULL,
    amount DECIMAL(5,2) DEFAULT 10.00,
    expired_at DATETIME NOT NULL DEFAULT (NOW() + INTERVAL 3 DAY),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS order_status (
    id INT NOT NULL PRIMARY KEY, 
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS order_data (
    id INT NOT NULL PRIMARY KEY,
    user_id INT NOT NULL,
    status_id INT NOT NULL DEFAULT 1,
    coupon_id INT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    CONSTRAINT `fk_order_user` FOREIGN KEY(user_id) REFERENCES user_info(id),
    CONSTRAINT `fk_order_status` FOREIGN KEY(status_id) REFERENCES order_status(id),
    CONSTRAINT `fk_order_coupon` FOREIGN KEY(coupon_id) REFERENCES coupon(id)
);

CREATE TABLE IF NOT EXISTS order_detail (
    id INT NOT NULL PRIMARY KEY, 
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    qty INT DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
    CONSTRAINT `fk_detail_order` FOREIGN KEY(order_id) REFERENCES order_data(id),
    CONSTRAINT `fk_detail_product` FOREIGN KEY(product_id) REFERENCES product(id) 
);