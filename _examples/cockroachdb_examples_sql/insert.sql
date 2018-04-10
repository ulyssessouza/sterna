-- cockroach sql --insecure --host 192.168.99.100 -p 32186 --database test < insert.sql
insert into test_tb (name) values('Ulysses');
