CREATE TABLE user (
  id int primary key auto_increment,
  first_name varchar(50),
  last_name varchar(50),
  email varchar(50),
  gender varchar(30),
  ip_address varchar(15),
  created_at datetime
);
