--
-- PostgreSQL database dump
--

\restrict 4bp1sGlXbx1r4Ud7d7BjUJ0YrVUi5uSys4crDR1cX9B72bQstuVqGrgS19tYUVc

-- Dumped from database version 16.10
-- Dumped by pg_dump version 16.10

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: cart_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cart_items (
                                   id integer NOT NULL,
                                   cart_id integer,
                                   product_id integer,
                                   size character varying(10) NOT NULL,
                                   quantity integer DEFAULT 1 NOT NULL,
                                   added_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.cart_items OWNER TO postgres;

--
-- Name: cart_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.cart_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.cart_items_id_seq OWNER TO postgres;

--
-- Name: cart_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.cart_items_id_seq OWNED BY public.cart_items.id;


--
-- Name: carts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.carts (
                              id integer NOT NULL,
                              user_id integer,
                              created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.carts OWNER TO postgres;

--
-- Name: carts_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.carts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.carts_id_seq OWNER TO postgres;

--
-- Name: carts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.carts_id_seq OWNED BY public.carts.id;


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.notifications (
                                      id integer NOT NULL,
                                      user_id integer,
                                      notification_type character varying(50) NOT NULL,
                                      message text NOT NULL,
                                      is_global boolean DEFAULT false,
                                      created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.notifications OWNER TO postgres;

--
-- Name: notifications_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.notifications_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.notifications_id_seq OWNER TO postgres;

--
-- Name: notifications_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.notifications_id_seq OWNED BY public.notifications.id;


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_items (
                                    id integer NOT NULL,
                                    order_id integer,
                                    product_id integer,
                                    product_name character varying(255) NOT NULL,
                                    size character varying(10) NOT NULL,
                                    quantity integer NOT NULL,
                                    price numeric(10,2) NOT NULL,
                                    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.order_items OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.order_items_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.order_items_id_seq OWNER TO postgres;

--
-- Name: order_items_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.order_items_id_seq OWNED BY public.order_items.id;


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
                               id integer NOT NULL,
                               user_id integer,
                               total_amount numeric(10,2) NOT NULL,
                               status character varying(50) DEFAULT 'PAID'::character varying,
                               delivery_address text NOT NULL,
                               phone_number character varying(20) NOT NULL,
                               payment_method character varying(50) DEFAULT 'CARD'::character varying,
                               created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.orders_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.orders_id_seq OWNER TO postgres;

--
-- Name: orders_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.orders_id_seq OWNED BY public.orders.id;


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
                                 id integer NOT NULL,
                                 name character varying(255) NOT NULL,
                                 description text,
                                 price numeric(10,2) NOT NULL,
                                 image_url character varying(500),
                                 created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                 updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                 highlight_text character varying(100),
                                 highlight_enabled boolean DEFAULT false
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.products_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.products_id_seq OWNER TO postgres;

--
-- Name: products_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.products_id_seq OWNED BY public.products.id;


--
-- Name: saved_cards; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.saved_cards (
                                    id integer NOT NULL,
                                    user_id integer,
                                    cardholder_name character varying(255) NOT NULL,
                                    card_number character varying(16) NOT NULL,
                                    expiry_month character varying(2) NOT NULL,
                                    expiry_year character varying(2) NOT NULL,
                                    cvv character varying(3) NOT NULL,
                                    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
                                    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.saved_cards OWNER TO postgres;

--
-- Name: saved_cards_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.saved_cards_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.saved_cards_id_seq OWNER TO postgres;

--
-- Name: saved_cards_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.saved_cards_id_seq OWNED BY public.saved_cards.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
                              id integer NOT NULL,
                              name character varying(255) NOT NULL,
                              email character varying(255) NOT NULL,
                              password_hash character varying(255) NOT NULL,
                              role character varying(50) DEFAULT 'USER'::character varying,
                              created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO postgres;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: cart_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items ALTER COLUMN id SET DEFAULT nextval('public.cart_items_id_seq'::regclass);


--
-- Name: carts id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts ALTER COLUMN id SET DEFAULT nextval('public.carts_id_seq'::regclass);


--
-- Name: notifications id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications ALTER COLUMN id SET DEFAULT nextval('public.notifications_id_seq'::regclass);


--
-- Name: order_items id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items ALTER COLUMN id SET DEFAULT nextval('public.order_items_id_seq'::regclass);


--
-- Name: orders id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders ALTER COLUMN id SET DEFAULT nextval('public.orders_id_seq'::regclass);


--
-- Name: products id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products ALTER COLUMN id SET DEFAULT nextval('public.products_id_seq'::regclass);


--
-- Name: saved_cards id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_cards ALTER COLUMN id SET DEFAULT nextval('public.saved_cards_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: cart_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cart_items (id, cart_id, product_id, size, quantity, added_at) FROM stdin;
3	1	2	S	1	2026-02-06 16:44:22.955761
4	1	7	S	1	2026-02-06 16:44:31.305951
\.


--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.carts (id, user_id, created_at) FROM stdin;
1	4	2026-02-02 22:19:48.553091
2	5	2026-02-02 22:34:26.098862
3	9	2026-02-11 16:55:10.306114
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.notifications (id, user_id, notification_type, message, is_global, created_at) FROM stdin;
1	4	order	Your order #1 has been successfully placed and paid! Total: $91.98	f	2026-02-03 20:38:55.317322
2	4	order	Your order #2 has been successfully placed and paid! Total: $45.99	f	2026-02-04 18:09:25.126718
3	9	order	Your order #3 has been successfully placed and paid! Total: $46.99	f	2026-02-11 19:08:06.12065
4	\N	promo	dasda	t	2026-02-11 19:19:18.845664
5	\N	promo	asdasd	t	2026-02-11 19:29:01.943605
\.


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_items (id, order_id, product_id, product_name, size, quantity, price, created_at) FROM stdin;
1	1	1	University Hoodie Blue	S	2	45.99	2026-02-03 20:38:55.317322
2	2	1	University Hoodie Blue	M	1	45.99	2026-02-04 18:09:25.126718
3	3	1	University Hoodie Blue	L	1	46.99	2026-02-11 19:08:06.12065
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, user_id, total_amount, status, delivery_address, phone_number, payment_method, created_at) FROM stdin;
1	4	91.98	PAID	Kabanbay batyr 62	+ 7777 777 77 77	CARD	2026-02-03 20:38:55.317322
2	4	45.99	PAID	Kabanbay batyr 62	+ 7777 777 77 77	CARD	2026-02-04 18:09:25.126718
3	9	46.99	PAID	Kabanbay batyr 62	+77777777777	CARD	2026-02-11 19:08:06.12065
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, description, price, image_url, created_at, updated_at, highlight_text, highlight_enabled) FROM stdin;
1	University Hoodie Blue	Premium quality hoodie with university logo. Made from 80% cotton, 20% polyester blend for maximum comfort and durability.	46.99	/static/images/University_Hoodie.png	2026-02-02 21:51:02.715735	2026-02-11 16:15:53.207625	Limited Drop	t
2	Classic T-Shirt White	Comfortable classic t-shirt featuring university emblem. Perfect for everyday wear, 100% cotton fabric.	19.99	/static/images/Classic_TShirt.png	2026-02-02 21:51:02.715735	2026-02-02 21:51:02.715735	\N	f
3	Sports Cap Red	Adjustable sports cap with embroidered university logo. One size fits all, breathable mesh design.	15.99	/static/images/Sports_Cap.png	2026-02-02 21:51:02.715735	2026-02-02 21:51:02.715735	\N	f
4	Varsity Jacket Black	Traditional varsity jacket with leather sleeves and university patches. Show your school pride in style.	89.99	/static/images/Varsity_Jacket.png	2026-02-02 21:51:02.715735	2026-02-02 21:51:02.715735	Trending Now	t
5	Sweatshirt Gray	Cozy sweatshirt with university name printed across chest. Warm fleece interior, perfect for cold days.	39.99	/static/images/Sweatshirt_Gray.png	2026-02-02 21:51:02.715735	2026-02-02 21:51:02.715735	\N	f
6	Track Pants Navy	Athletic track pants with university stripes. Comfortable elastic waistband and zippered pockets.	34.99	/static/images/Track_Pants.png	2026-02-02 21:51:02.715735	2026-02-02 21:51:02.715735	\N	f
7	Zip Hoodie Green	Full-zip hoodie with front pockets and university crest. Premium quality with reinforced stitching.	49.99	/static/images/Zip_Hoodie.png	2026-02-02 21:51:02.715735	2026-02-02 21:51:02.715735	\N	f
8	Baseball Jersey White	Official university baseball jersey with team colors. Breathable mesh fabric, perfect for sports or casual wear.	54.99	/static/images/Baseball_Jersey.png	2026-02-02 21:51:02.715735	2026-02-02 21:51:02.715735	\N	f
\.


--
-- Data for Name: saved_cards; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.saved_cards (id, user_id, cardholder_name, card_number, expiry_month, expiry_year, cvv, created_at, updated_at) FROM stdin;
1	4	Bek G	4539123456789012	07	32	444	2026-02-03 20:38:55.317322	2026-02-03 20:38:55.317322
2	9	Zhan	1234412365238520	10	36	456	2026-02-11 19:08:06.12065	2026-02-11 19:08:06.12065
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, email, password_hash, role, created_at) FROM stdin;
1	Test User	test@university.edu	$2a$10$N9qo8uLOickgx2ZMRZoMye7FRNv6qx6vUwPjgmJmJSMl0LGJGz5Iq	USER	2026-02-02 21:51:22.704492
4	Bekbauly	bekbaulykalymzan@gmail.com	$2a$10$GFI.HY.dhUbyNlEauwlVCuIkD4sOAajJTw8wWFUcqDcE/522XLJk.	USER	2026-02-02 22:19:34.195284
5	Admin	admin@university.edu	$2a$10$HB.fN.caYkIxgeT7s5169.2nm7qoCq0Brg1qW8x0DTAI.GnUI5pzC	ADMIN	2026-02-02 22:33:03.916238
7	Koke	Kazakh@gmail.com	$2a$10$g9GOJaYBmAbsKL69YxIRxuGxr8WBsKeLyEA.TNDEQfG7HJD3yWJxG	USER	2026-02-03 19:33:43.995308
8	Koke	ggg@gmail.com	$2a$10$0VtNbGCe/Nm96/2KnbqbIeVCGORClvS8lLDiFU.AlTU8KtpBnQgVq	USER	2026-02-06 16:47:24.108435
9	Zhan	Zhan@gmail.com	$2a$10$fpMJM6ZVj0tb9g.AgsgNgOsTQ0VA0ICuX39OOga0MUq8XwAO3TyUK	USER	2026-02-11 16:54:46.406752
\.


--
-- Name: cart_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.cart_items_id_seq', 5, true);


--
-- Name: carts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.carts_id_seq', 3, true);


--
-- Name: notifications_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.notifications_id_seq', 5, true);


--
-- Name: order_items_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.order_items_id_seq', 3, true);


--
-- Name: orders_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.orders_id_seq', 3, true);


--
-- Name: products_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.products_id_seq', 11, true);


--
-- Name: saved_cards_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.saved_cards_id_seq', 2, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.users_id_seq', 9, true);


--
-- Name: cart_items cart_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_pkey PRIMARY KEY (id);


--
-- Name: carts carts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_pkey PRIMARY KEY (id);


--
-- Name: carts carts_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_user_id_key UNIQUE (user_id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: saved_cards saved_cards_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_cards
    ADD CONSTRAINT saved_cards_pkey PRIMARY KEY (id);


--
-- Name: saved_cards saved_cards_user_id_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_cards
    ADD CONSTRAINT saved_cards_user_id_key UNIQUE (user_id);


--
-- Name: cart_items unique_cart_product_size; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT unique_cart_product_size UNIQUE (cart_id, product_id, size);


--
-- Name: users users_email_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: idx_cart_items_cart_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_cart_items_cart_id ON public.cart_items USING btree (cart_id);


--
-- Name: idx_notifications_global; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_notifications_global ON public.notifications USING btree (is_global);


--
-- Name: idx_notifications_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_notifications_user_id ON public.notifications USING btree (user_id);


--
-- Name: idx_order_items_order_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_order_items_order_id ON public.order_items USING btree (order_id);


--
-- Name: idx_orders_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_orders_user_id ON public.orders USING btree (user_id);


--
-- Name: idx_saved_cards_user_id; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_saved_cards_user_id ON public.saved_cards USING btree (user_id);


--
-- Name: idx_users_role; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX idx_users_role ON public.users USING btree (role);


--
-- Name: cart_items cart_items_cart_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_cart_id_fkey FOREIGN KEY (cart_id) REFERENCES public.carts(id) ON DELETE CASCADE;


--
-- Name: cart_items cart_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: carts carts_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: notifications notifications_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_order_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_order_id_fkey FOREIGN KEY (order_id) REFERENCES public.orders(id) ON DELETE CASCADE;


--
-- Name: order_items order_items_product_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.products(id) ON DELETE CASCADE;


--
-- Name: orders orders_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: saved_cards saved_cards_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.saved_cards
    ADD CONSTRAINT saved_cards_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict 4bp1sGlXbx1r4Ud7d7BjUJ0YrVUi5uSys4crDR1cX9B72bQstuVqGrgS19tYUVc

