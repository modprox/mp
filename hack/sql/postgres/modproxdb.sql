create table modules (
  id serial primary key,
  source varchar(1024) not null,
  version varchar(128) not null,
  created timestamp not null default current_timestamp,
  unique(source, version)
);

create table proxy_configurations (
  id serial primary key,
  hostname varchar(128) not null,
  port integer not null,
  storage text not null,
  registry text not null,
  transforms text not null,
  ts timestamp not null default current_timestamp,
  unique(hostname, port)
);

create table proxy_heartbeats (
  id serial primary key,
  hostname varchar(128) not null,
  port integer not null,
  num_packages integer not null,
  num_modules integer not null,
  ts timestamp not null default current_timestamp,
  unique(hostname, port)
);
