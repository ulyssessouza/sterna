-- cockroach sql --insecure --host 192.168.99.100 -p 30257 --database test < insert.sql
insert into test_tb (name) values('Ulysses');
