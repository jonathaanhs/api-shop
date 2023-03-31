CREATE TYPE promo_type_enum as enum('product', 'discount');

CREATE TABLE products (
	product_id bigserial NOT NULL,
	sku varchar(255) NOT NULL,
	"name" varchar(255) NOT NULL,
	price numeric(50, 3) NOT NULL,
	qty int4 NOT NULL,
	CONSTRAINT product_id_pkey PRIMARY KEY (product_id)
);

CREATE TABLE promos (
	promo_id bigserial NOT NULL,
	product_id int8 NOT NULL,
	promo_type promo_type_enum NOT NULL,
	reward numeric(50, 3) NOT NULL,
	min_qty int4 NOT NULL,
	CONSTRAINT promo_id_pkey PRIMARY KEY (promo_id)
);

CREATE TABLE orders (
	order_id bigserial NOT NULL,
	"date" timestamp NOT NULL DEFAULT now(),
	total numeric(50, 3) NOT NULL,
	CONSTRAINT order_id_pkey PRIMARY KEY (order_id)
);

CREATE TABLE order_details (
	order_detail_id bigserial NOT NULL,
	order_id int8 NOT NULL,
	product_id int8 NOT NULL,
	promo_id int8 NOT NULL,
	price numeric(50, 3) NOT NULL,
	qty int4 NOT NULL,
	CONSTRAINT order_detail_id_pkey PRIMARY KEY (order_detail_id)
);


