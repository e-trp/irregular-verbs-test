create table if not exists verbs (
id INTEGER PRIMARY KEY AUTOINCREMENT,
infinitive VARCHAR(50) NOT NULL,
past_simple VARCHAR(50) NOT NULL,
past_participle VARCHAR(50) NOT NULL
);

create table if not exists  tests (
id INTEGER PRIMARY KEY AUTOINCREMENT, 
start_datetime VARCHAR(100)  NOT NULL,
end_datetime VARCHAR(100),
user_name VARCHAR(100) NOT NULL,
words_count INTEGER NOT NULL
);


create table if not exists  errors (
id INTEGER PRIMARY KEY AUTOINCREMENT,
test_id INTEGER NOT NULL,
source VARCHAR(50) NOT NULL,
user_guess VARCHAR(50),
FOREIGN KEY(test_id) REFERENCES artist(tests)
);
