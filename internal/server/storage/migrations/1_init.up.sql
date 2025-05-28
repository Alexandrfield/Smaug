CREATE TABLE if NOT EXISTS Users (id SERIAL PRIMARY KEY, login text, password bytea, signKeyComplicated bytea, key bytea);
CREATE TABLE if NOT EXISTS EncryptedData (id SERIAL PRIMARY KEY, login text, info text, data bytea);
