DROP TABLE IF EXISTS user;

CREATE TABLE user (
  id varbinary(16) primary key,
  name varchar(50),
  created_at datetime(6) default now(6)
)