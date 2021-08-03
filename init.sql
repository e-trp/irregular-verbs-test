create table if not exists verbs (
id INTEGER PRIMARY KEY AUTOINCREMENT,
infinit VARCHAR(50) NOT NULL,
simple VARCHAR(50) NOT NULL,
participl VARCHAR(50) NOT NULL
);

create table if not exists  tests (
id INTEGER PRIMARY KEY AUTOINCREMENT, 
datetime_stamp VARCHAR(100)  NOT NULL,
user_name NOT NULL
);


create table if not exists  errors (
id INTEGER PRIMARY KEY AUTOINCREMENT,
test_id INTEGER NOT NULL,
source VARCHAR(50) NOT NULL,
user_guess VARCHAR(50),
FOREIGN KEY(test_id) REFERENCES artist(tests)
);
