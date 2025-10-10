DELIMITER //
-- before insert check current stock, 
-- if not found show error
-- if found check if current quantity is enough if not enough error
-- if enough then update stock
-- #####################
-- after insert
-- update stock movement
-- ##################
CREATE TRIGGER od_before_ins
BEFORE INSERT ON order_details
FOR EACH ROW
BEGIN
  DECLARE cur_qty INT;

  SELECT qty_on_hand INTO cur_qty
  FROM stock
  WHERE product_id = NEW.product_id 
  FOR UPDATE;

  IF cur_qty IS NULL THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Stock row not found for product';
  END IF;

  IF cur_qty < NEW.quantity THEN
    SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = 'Insufficient stock';
  END IF;

  UPDATE stock
  SET qty_on_hand = qty_on_hand - NEW.quantity
  WHERE product_id = NEW.product_id;
END//
//

CREATE TRIGGER od_after_ins
AFTER INSERT ON order_details
FOR EACH ROW
BEGIN
  INSERT INTO stock_movements
    (product_id, order_detail_id, qty_change, reason)
  VALUES
    (NEW.product_id, NEW.order_detail_id, -NEW.quantity, 'order');
END//
//
DELIMITER ;
