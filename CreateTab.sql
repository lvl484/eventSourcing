CREATE TABLE event
(
  version bigserial PRIMARY KEY,
  id varchar NOT NULL,
  data jsonb NOT NULL,
  time varchar NOT NULL
);
