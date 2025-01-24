CREATE DATABASE rec_twitcasting;

-- connection
\c rec_twitcasting;

-- create sequence
CREATE SEQUENCE speakers_id_seq START 1;

-- create table
CREATE TABLE IF NOT EXISTS speakers (
  id INTEGER PRIMARY KEY DEFAULT nextval('speakers_id_seq'),
  username VARCHAR(255) NOT NULL,
  recording_state BOOLEAN NOT NULL DEFAULT FALSE,
  created_date_time timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- insert data (sample)
-- INSERT INTO speakers (username) VALUES('twitcasting_username');
