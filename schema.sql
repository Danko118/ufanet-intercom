--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
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
-- Name: events; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.events (
    id integer NOT NULL,
    mac character varying(17) NOT NULL,
    event_name text NOT NULL,
    event_args integer,
    event_desc text,
    event_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.events OWNER TO postgres;

--
-- Name: events_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.events_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.events_id_seq OWNER TO postgres;

--
-- Name: events_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.events_id_seq OWNED BY public.events.id;


--
-- Name: intercoms; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.intercoms (
    id integer NOT NULL,
    mac character varying(17) NOT NULL,
    address text,
    vendor text,
    aparts integer[]
);


ALTER TABLE public.intercoms OWNER TO postgres;

--
-- Name: intercoms_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.intercoms_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.intercoms_id_seq OWNER TO postgres;

--
-- Name: intercoms_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.intercoms_id_seq OWNED BY public.intercoms.id;


--
-- Name: intercomstatus; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.intercomstatus (
    mac character varying(17),
    status text NOT NULL,
    "time" timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT intercomstatus_status_check CHECK ((status = ANY (ARRAY['healthy'::text, 'broken'::text, 'unreachable'::text, 'resolving'::text])))
);


ALTER TABLE public.intercomstatus OWNER TO postgres;

--
-- Name: events id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events ALTER COLUMN id SET DEFAULT nextval('public.events_id_seq'::regclass);


--
-- Name: intercoms id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercoms ALTER COLUMN id SET DEFAULT nextval('public.intercoms_id_seq'::regclass);


--
-- Name: events events_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_pkey PRIMARY KEY (id);


--
-- Name: intercoms intercoms_mac_key; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercoms
    ADD CONSTRAINT intercoms_mac_key UNIQUE (mac);


--
-- Name: intercoms intercoms_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercoms
    ADD CONSTRAINT intercoms_pkey PRIMARY KEY (id);


--
-- Name: events events_mac_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.events
    ADD CONSTRAINT events_mac_fkey FOREIGN KEY (mac) REFERENCES public.intercoms(mac);


--
-- Name: intercomstatus intercomstatus_mac_fkey; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.intercomstatus
    ADD CONSTRAINT intercomstatus_mac_fkey FOREIGN KEY (mac) REFERENCES public.intercoms(mac) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

