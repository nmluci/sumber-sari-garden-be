CREATE TRIGGER `trig_add_trx` BEFORE UPDATE ON order_data 
    FOR EACH ROW 
    BEGIN
        DECLARE p_id INT;
        DECLARE p_qty INT;
        DECLARE product_qty INT;

        SET p_id = (SELECT product_id FROM order_detail WHERE order_id = NEW.id);
		SET p_qty = (SELECT qty FROM order_detail WHERE order_id = NEW.id);
        SET product_qty = (SELECT qty FROM product WHERE id=p_id);

        IF (p_qty >= product_qty) THEN
            SET p_qty = product_qty;
            UPDATE order_detail d SET d.qty = product_qty WHERE order_id=NEW.id;
        END IF;
	        
        IF (NEW.status_id >= 2 ) THEN
            IF (p_qty <= 0) THEN
                UPDATE product p SET p.qty = 0 WHERE p.id=p_id;
            ELSE
                UPDATE product p SET p.qty = (p.qty - p_qty) WHERE p.id = p_id;
            END IF;
        END IF;
    END;

CREATE TRIGGER `trig_delete_category` BEFORE DELETE ON product_category
    FOR EACH ROW
    BEGIN
        UPDATE product p SET p.category_id=1 WHERE p.category_id=OLD.id;
    END;