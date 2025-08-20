-- public.m_item definition

-- Drop table

-- DROP TABLE public.m_item;

CREATE TABLE public.m_item (
	id bigserial NOT NULL,
	m_bu_id int8 NULL,
	code varchar(50) NULL,
	m_item_type_id int8 NULL,
	m_cat1_id int8 NULL,
	m_cat2_id int8 NULL,
	m_cat3_id int8 NULL,
	m_cat4_id int8 NULL,
	item_name varchar(100) NOT NULL,
	item_name_long text NULL,
	unit_id int8 NULL,
	unit varchar(100) NULL,
	mnfct varchar(100) NULL,
	price_base numeric(18, 2) NOT NULL,
	item_photo varchar(255) NULL,
	spec varchar(100) NULL,
	weight numeric(8, 2) NULL,
	weight_unit_id int8 NULL,
	dim_l float8 NULL,
	dim_l_unit_id int8 NULL,
	dim_p float8 NULL,
	dim_p_unit_id int8 NULL,
	dim_t float8 NULL,
	dim_t_unit_id int8 NULL,
	is_active bool DEFAULT true NOT NULL,
	creator_id int4 NULL,
	editor_id int4 NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	is_timbangan bool DEFAULT false NULL,
	round numeric(10, 2) NULL,
	flag_ppn bool DEFAULT false NULL,
	m_supp_id int8 NULL,
	default_price_sale numeric(18, 2) NULL,
	barcode varchar(255) NULL,
	CONSTRAINT m_item_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_m_item_code ON public.m_item USING btree (code);
CREATE INDEX idx_m_item_id ON public.m_item USING btree (id);
CREATE INDEX idx_m_item_unit_id ON public.m_item USING btree (unit_id);