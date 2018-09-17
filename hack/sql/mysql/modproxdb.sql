create table modules (
  id int(3) unsigned not null auto_increment,
  source varchar(896) not null,
  version varchar(128) not null,
  created timestamp not null default current_timestamp,
  primary key(id),
  unique (source, version)
) engine=InnoDB default charset=utf8;

create table proxy_configurations (
  id int(3) unsigned not null auto_increment,
  hostname varchar(128) not null,
  port int(6) not null,
  storage text not null,
  registry text not null,
  transforms text not null,
  ts timestamp not null default current_timestamp,
  primary key(id),
  unique (hostname, port)
) engine=InnoDB default charset=utf8;

create table proxy_heartbeats (
  id int(3) unsigned not null auto_increment,
  hostname varchar(128) not null,
  port int(6) not null,
  num_modules int(10) not null,
  num_versions int(10) not null,
  ts timestamp not null default current_timestamp,
  primary key(id),
  unique (hostname, port)
) engine=InnoDB default charset=utf8;
