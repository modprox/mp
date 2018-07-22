create database if not exists modproxdb;

use modproxdb;

create table registry (
  id int(5) unsigned not null auto_increment,
  source varchar(924) not null,
  version varchar(100) not null,
  created timestamp not null default current_timestamp,
  primary key(id),
  unique(source, version)
) engine=InnoDB default charset=utf8;
