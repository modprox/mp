create table proxy_module_zips (
  id int(5) unsigned not null auto_increment,
  path varchar(1024) not null,
  zip blob not null,
  primary key(id),
  unique (path)
) engine=InnoDB default charset=utf8;

create table proxy_modules_index (
  id int(5) unsigned not null auto_increment,
  source varchar(256) not null,
  version varchar(256) not null,
  registry_mod_id int(3) unsigned not null, # references modules.id, which may be in a different db
  go_mod_contents blob not null,
  rev_info_contents blob not null,
  primary key(id),
  unique (source, version)
) engine=InnoDB default charset=utf8;
