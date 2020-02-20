drop table if exists users;
create table users (
  user_id integer primary key autoincrement,
  username string not null,
  email string not null,
  pw_hash string not null,
  image_url string
);

drop table if exists followers;
create table followers (
  who_id integer,
  whom_id integer
);

drop table if exists messages;
create table messages (
  message_id integer primary key autoincrement,
  author_id integer not null,
  text string not null,
  pub_date integer,
  flagged integer
);
