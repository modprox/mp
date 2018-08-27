create database if not exists modproxdb;

use modproxdb;

create table sources (
  id int(5) unsigned not null auto_increment,
  source varchar(1024) not null,
  created timestamp not null default current_timestamp,
  primary key(id),
  unique(source)
) engine=InnoDB default charset=utf8;

create table tags (
  id int(5) unsigned not null auto_increment,
  tag varchar(128) not null,
  created timestamp not null default current_timestamp,
  source_id int(5) unsigned not null,
  primary key(id),
  foreign key (source_id) references sources(id) on delete cascade
) engine=InnoDB default charset=utf8;

create table redirects (
  id int(3) unsigned not null auto_increment,
  original varchar(128) not null,
  substitution varchar(128) not null,
  created timestamp not null default current_timestamp,
  primary key(id),
  unique(original)
) engine=InnoDB default charset=utf8;
