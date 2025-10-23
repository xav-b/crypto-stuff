/*
 * Initialise PostgreSQL
 */

 CREATE DATABASE crypto;

 \connect crypto;

 CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
 CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;
