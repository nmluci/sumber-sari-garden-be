ALTER TABLE coupon
    ADD COLUMN description VARCHAR(255);

UPDATE coupon 
SET description = "GRATIS SOB"
WHERE id = 1;

UPDATE coupon 
SET description = "DISKON NICH"
WHERE id = 2;

UPDATE coupon 
SET description = "SUMBER REJEKI"
WHERE id = 3;