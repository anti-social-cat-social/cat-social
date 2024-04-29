create table users (
 uuid uuid PRIMARY KEY,
 name varchar(128) not null,
 email varchar(64) unique not null,
 password varchar(255) not null,
 created_at timestamp default current_timestamp
)
