create table modules (
 -- todo: go back and do all this for mysql
);

-- upsert on proxy startup, so registry viewers can
-- see what the file configuration was for each proxy
-- that is connected to the registry
create table proxy_configurations (
  id int(3) unsigned not null auto_increment,
  hostname varchar(128) not null,
  port int(6) not null,
  transforms text not null,
  ts timestamp not null default current_timestamp,
  primary key(id),
  unique (hostname, port)
) engine=InnoDB default charset=utf8;

-- upsert on a polling period, so registry viewers
-- can see what each proxy has been up to recently.
create table proxy_heartbeats (
  id int(3) unsigned not null auto_increment,
  hostname varchar(128) not null,
  port int(6) not null,
  num_packages int(10) not null,
  num_modules int(10) not null,
  ts timestamp not null default current_timestamp,
  primary key(id),
  unique (hostname, port)
) engine=InnoDB default charset=utf8;
