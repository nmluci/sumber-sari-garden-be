CREATE TRIGGER `trig_add_trx` AFTER UPDATE ON order_data 
    FOR EACH ROW 
    BEGIN
        DECLARE p_id INT;
        DECLARE p_qty INT;
            
        IF (NEW.status_id = 2 ) THEN
            SELECT product_id, qty INTO p_qty, p_id FROM order_detail WHERE order_id = NEW.id;
            UPDATE product p SET p.qty = p.qty - p_qty WHERE p.id = p_id;
            END IF;
    END;