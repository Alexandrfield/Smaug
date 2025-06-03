CREATE TABLE if NOT EXISTS Users (id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY, login text, password bytea, signKeyComplicated bytea, key bytea);
CREATE TABLE if NOT EXISTS EncryptedData (id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY, login text, info text, data bytea);
