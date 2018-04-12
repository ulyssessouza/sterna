CREATE TABLE test_tb (
	"name" STRING NULL,
	FAMILY "primary" ("name", rowid)
);

INSERT INTO test_tb ("name") VALUES
	('Ulysses');
