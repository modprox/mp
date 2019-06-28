create table proxy_module_zips (
  id int(5) unsigned not null auto_increment,
  s_at_v varchar(1024) not null, -- unique module identifier source@version
  zip mediumblob not null, -- binary blob of the well formed zip archive
  primary key(id),
  unique (s_at_v)
) engine=InnoDB default charset=utf8;

create table proxy_modules_index (
  id int(5) unsigned not null auto_increment,
  source varchar(256) not null, -- module package, e.g. github.com/pkg/errors
  version varchar(256) not null, -- module version, e.g. v1.0.0-alpha1
  registry_mod_id int(5) unsigned not null, -- registry serial number of the module
  go_mod_file text not null, -- text of the go.mod file of the module
  version_info text not null, -- JSON of .info pseudo file
  primary key(id),
  unique (source, version)
) engine=InnoDB default charset=utf8;
