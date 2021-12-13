CREATE TABLE IF NOT EXISTS public.product_stock
(
    deposit character varying(12) COLLATE pg_catalog."default" NOT NULL,
    location character varying(4) COLLATE pg_catalog."default" NOT NULL,
    product_id character varying(60) COLLATE pg_catalog."default" NOT NULL,
    quantity integer,
    CONSTRAINT product_stock_pkey PRIMARY KEY (deposit, location, product_id)
);