-- public.m_supp definition

-- Drop table

-- DROP TABLE public.m_supp;

CREATE TABLE public.m_supp (
	id bigserial NOT NULL,
	code varchar(255) NULL,
	origin varchar(100) NULL,
	group_supp_id int8 NULL,
	"type" varchar(100) NULL,
	"name" varchar(100) NOT NULL,
	nib varchar(100) NULL,
	top_id int8 NULL,
	flag_ppn bool DEFAULT false NOT NULL,
	npwp varchar(255) NULL,
	addr text NULL,
	prov_id int8 NULL,
	city_id int8 NULL,
	district_id int8 NULL,
	post_code varchar(10) NULL,
	coa_hutang_id int8 NULL,
	phone1 varchar(100) NULL,
	phone2 varchar(100) NULL,
	cp1 varchar(255) NULL,
	cp1_phone varchar(100) NULL,
	cp2 varchar(255) NULL,
	cp2_phone varchar(100) NULL,
	"desc" text NULL,
	is_active bool DEFAULT true NOT NULL,
	creator_id int4 NULL,
	editor_id int4 NULL,
	created_at timestamp(0) NULL,
	updated_at timestamp(0) NULL,
	CONSTRAINT m_supp_pkey PRIMARY KEY (id)
);
CREATE INDEX idx_m_supp_id ON public.m_supp USING btree (id);
CREATE INDEX idx_m_supp_id_code ON public.m_supp USING btree (id, code);