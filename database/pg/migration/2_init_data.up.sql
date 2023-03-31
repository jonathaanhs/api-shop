INSERT INTO products (sku,"name",price,qty) VALUES
	 ('120P90','Google Home',49.990,10),
	 ('43N23P','MacBook Pro',5399.990,5),
	 ('A304SD','Alexa Speaker',109.500,10),
	 ('234234','Raspberry Pi B',30.000,2);

INSERT INTO promos (product_id,promo_type,reward,min_qty) VALUES
	 (1,'product'::promo_type_enum,1.000,3),
	 (2,'product'::promo_type_enum,4.000,1),
	 (3,'discount'::promo_type_enum,10.000,4);