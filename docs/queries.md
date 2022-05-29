# Query List

- [Query List](#query-list)
    - [Store User Info](#store-user-info)
    - [Store User Auth Info](#store-user-auth-info)
    - [Get User by Email](#get-user-by-email)
    - [Get User by ID](#get-user-by-id)
    - [Get User Info by ID](#get-user-info-by-id)
    - [Count Products](#count-products)
    - [Search Product by Keyword With Limit and Offset](#search-product-by-keyword-with-limit-and-offset)
    - [Get Product by ID](#get-product-by-id)
    - [Insert Product](#insert-product)
    - [Update Product Data](#update-product-data)
    - [Delete Product by ID](#delete-product-by-id)
    - [Get Product Categories](#get-product-categories)
    - [Get Product Category by ID](#get-product-category-by-id)
    - [Store Product Category](#store-product-category)
    - [Update Product Category](#update-product-category)
    - [Delete Product Category](#delete-product-category)
    - [Get Active Coupons](#get-active-coupons)
    - [Get Order Details by ID](#get-order-details-by-id)
    - [Get Order Details by ID](#get-order-details-by-id-1)
    - [Get Order Metadata](#get-order-metadata)
    - [Get Order Metadata All](#get-order-metadata-all)
    - [Get Order Data](#get-order-data)
    - [Start New Order](#start-new-order)
    - [Get Order's Item Metadata](#get-orders-item-metadata)
    - [Insert New Item to Order](#insert-new-item-to-order)
    - [Update Order Data](#update-order-data)
    - [Delete Order Data by ID](#delete-order-data-by-id)
    - [Get Coupon by Code](#get-coupon-by-code)
    - [Checkout Order](#checkout-order)
    - [Get Unpaid Order](#get-unpaid-order)
    - [Verify Order](#verify-order)

### Store User Info
```sql
INSERT INTO user_info(user_id, first_name, last_name, phone, address) VALUES (?, ?, ?, ?, ?)
```

### Store User Auth Info
```sql
INSERT INTO user (role_id, email, password) VALUES (?, ?, ?);
```

### Get User by Email
```sql
SELECT id, email, password, role_id FROM user WHERE email=?
```

### Get User by ID
```sql
SELECT id, email, password, role_id FROM user WHERE id=?
```

### Get User Info by ID
```sql
SELECT first_name, last_name, phone, address FROM user_info WHERE id=?;
```

### Count Products
```sql
SELECT COUNT(*) FROM product
```

### Search Product by Keyword With Limit and Offset
```sql
SELECT p.id, p.category_id, p.name, pc.name, p.price, p.qty, p.url, p.description FROM product p 
    LEFT JOIN product_category pc ON p.category_id=pc.id 
    WHERE p.name LIKE "%%%s%%" AND pc.id IN (%s) ORDER BY %s LIMIT ? OFFSET
```

### Get Product by ID
```sql
SELECT p.id, p.category_id, p.name, p.price, p.qty, p.url, p.description FROM product p WHERE p.id=?
```

### Insert Product
```sql

INSERT INTO product(category_id, name, price, qty, url, description) VALUES (?, ?, ?, ?, ?, ?)
```

### Update Product Data
```sql
UPDATE product SET category_id = ?, name=?, price=?, qty=?, url=?, description=? WHERE id=?
```

### Delete Product by ID
```sql
DELETE FROM product WHERE id=?
```

### Get Product Categories
```sql
SELECT c.id, c.name FROM product_category c
```

### Get Product Category by ID
```sql
SELECT c.id, c.name FROM product_category c WHERE c.id=?
```

### Store Product Category
```sql
INSERT INTO product_category(name) VALUES (?)
```

### Update Product Category
```sql
UPDATE product_category SET name=? WHERE id=?
```

### Delete Product Category
```sql
DELETE FROM product_category WHERE id=?
```

### Get Active Coupons
```sql
SELECT c.id, c.code, c.amount, c.description, c.expired_at FROM coupon c 
    WHERE c.expired_at > NOW() LIMIT ? OFFSET ?
```

### Get Order Details by ID
```sql
SELECT od.id, od.order_id, od.product_id, p.name, p.price, od.qty, (
    case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end) disc, 
    ((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) subtotal 
    FROM order_detail od LEFT JOIN product p ON od.product_id=p.id LEFT JOIN order_data o ON od.order_id=o.id 
    LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.id=?
```

### Get Order Details by ID
```sql
SELECT o.id, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total,  
    COUNT(*) item_count 
    FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
    LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN coupon c ON o.coupon_id=c.id 
    GROUP BY o.id HAVING o.id=?
```

### Get Order Metadata
```sql

SELECT * FROM (
        SELECT o.id, o.user_id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
        COUNT(*) item_count, c.code, s.name 
        FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
        LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id 
        LEFT JOIN coupon c ON o.coupon_id=c.id 
        WHERE o.status_id = ? GROUP BY o.id, o.user_id HAVING o.user_id=? 
    UNION 
        SELECT o.id, o.user_id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
        COUNT(*) item_count, c.code, s.name 
        FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
        LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id LEFT JOIN coupon c ON o.coupon_id=c.id 
        WHERE o.status_id = ? GROUP BY o.id, o.user_id HAVING o.user_id=? 
    UNION 
        SELECT o.id, o.user_id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
        COUNT(*) item_count, c.code, s.name 
        FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
        LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id LEFT JOIN coupon c ON o.coupon_id=c.id 
        WHERE o.status_id = ? GROUP BY o.id, o.user_id HAVING o.user_id=?
    ) t WHERE (t.created_at BETWEEN ? AND ?) LIMIT ? OFFSET ? 
```

### Get Order Metadata All
```sql
SELECT * FROM (
        SELECT o.id, o.user_id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
        COUNT(*) item_count, c.code, s.name 
        FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
        LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id 
        LEFT JOIN coupon c ON o.coupon_id=c.id WHERE o.status_id = ? GROUP BY o.id, o.user_id 
    UNION 
        SELECT o.id, o.user_id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
        COUNT(*) item_count, c.code, s.name 
        FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
        LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id, LEFT JOIN coupon c ON o.coupon_id=c.id 
        WHERE o.status_id = ? GROUP BY o.id, o.user_id 
    UNION 
        SELECT o.id, o.user_id, o.created_at, SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
        COUNT(*) item_count, c.code, s.name 
        FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
        LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN order_status s ON s.id=o.status_id 
        LEFT JOIN coupon c ON o.coupon_id=c.id 
        WHERE o.status_id = ? GROUP BY o.id, o.user_id) t WHERE (t.created_at BETWEEN ? AND ?
    ) LIMIT ? OFFSET ? | Metadata All
```

### Get Order Data
```sql
SELECT o.id, o.user_id, o.status_id, s.name, o.coupon_id FROM order_data o 
    LEFT JOIN order_status s ON o.status_id=s.id 
    WHERE o.user_id=? AND o.status_id=1 ORDER BY o.created_at DESC LIMIT 1
```

### Start New Order
```sql
INSERT INTO order_data(user_id) VALUES (?)
```

### Get Order's Item Metadata 
```sql
SELECT od.id, od.order_id, od.product_id, p.name, p.price, od.qty, 
    (case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end) disc, 
    ((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) subtotal 
    FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
    LEFT JOIN order_data o ON od.order_id=o.id LEFT JOIN coupon c ON o.coupon_id=c.id 
    WHERE o.id=? AND p.id = ?
```

### Insert New Item to Order
```sql
INSERT INTO order_detail(order_id, product_id, qty) VALUES (?, ?, ?)
```

### Update Order Data
```sql
UPDATE order_detail SET qty=? WHERE product_id=? AND order_id=?
```

### Delete Order Data by ID
```sql
DELETE FROM order_detail WHERE product_id=? AND order_id=?
```

### Get Coupon by Code
```sql
SELECT c.id, c.code, (c.amount/100), c.expired_at FROM coupon c 
    WHERE c.code=? AND c.expired_at > NOW()
```

### Checkout Order
``` sql
UPDATE order_data SET created_at=NOW(), status_id=2, coupon_id=? 
    WHERE order_data.user_id=? AND order_data.id=?
```

### Get Unpaid Order
```sql
SELECT o.id, 
    SUM((p.price*od.qty)-(case when o.coupon_id is not null then ((c.amount/100)*(p.price*od.qty)) else 0 end)) grand_total, 
    COUNT(*) item_count 
    FROM order_detail od LEFT JOIN product p ON od.product_id=p.id 
    LEFT JOIN order_data o ON od.order_id=o.id 
    LEFT JOIN coupon c ON o.coupon_id=c.id 
    WHERE o.status_id=2 GROUP BY o.id
```

### Verify Order
```sql
UPDATE order_data SET status_id=3 WHERE order_data.id=?
```