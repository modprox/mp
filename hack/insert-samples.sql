/*
To run this, launch mysql client (using connect-mysql.sh) and run the command:
source insert-samples.sql
*/

insert into sources (source) values ("github.com/foo/bar");
insert into sources (source) values ("github.com/foo/baz");
insert into sources (source) values ("github.com/other/a");
insert into sources (source) values ("github.com/other/b");
insert into sources (source) values ("github.com/other/c");
insert into sources (source) values ("github.com/other/d"); /* unused */

insert into tags (tag, source_id) values ("v1.0.0", 1);
insert into tags (tag, source_id) values ("v1.1.0", 1);
insert into tags (tag, source_id) values ("v0.0.1", 2);
insert into tags (tag, source_id) values ("v0.0.2", 1);
insert into tags (tag, source_id) values ("v2.0.0", 3);
insert into tags (tag, source_id) values ("v1.3.1", 1);
insert into tags (tag, source_id) values ("v0.1.0", 2);
insert into tags (tag, source_id) values ("v0.0.2", 2);
insert into tags (tag, source_id) values ("v1.1.1", 1);
insert into tags (tag, source_id) values ("v3.0.0", 3);
insert into tags (tag, source_id) values ("v1.8.0", 3);
insert into tags (tag, source_id) values ("v1.2.1", 4);
insert into tags (tag, source_id) values ("v2.3.0", 5);
insert into tags (tag, source_id) values ("v0.4.4", 4);

