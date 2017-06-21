--user
DROP TABLE IF EXISTS film_genre;
DROP TABLE IF EXISTS rented_film;
DROP TABLE IF EXISTS film;
DROP TABLE IF EXISTS genre;
DROP TABLE IF EXISTS person;
CREATE TABLE person (
  id serial primary key,
  username TEXT,
  password TEXT NOT NULL,
  login TEXT UNIQUE NOT NULL,
  age INTEGER,
  telephone TEXT,
  created_date timestamp
);

INSERT INTO person(username, password, login, age, telephone, created_date)
  VALUES('Anton Lempiy', '1fc854110e5532480000542834f453de31936c2f',
        'lempiy', 24, '+3806612345678', now());

--film
CREATE TABLE film (
  id serial primary key,
  name text not null,
  year INTEGER,
  added_at timestamp
);


--genre
CREATE TABLE genre (
  id serial primary key,
  name TEXT not null,
  added_at timestamp
);

--film_genre
CREATE TABLE film_genre (
  id serial primary key,
  film_id INTEGER REFERENCES film(id) ON UPDATE CASCADE ON DELETE CASCADE,
  genre_id INTEGER REFERENCES genre(id) ON UPDATE CASCADE ON DELETE CASCADE,
  added_at timestamp
);

INSERT INTO genre(name, added_at) VALUES('Comedy', now());
INSERT INTO genre(name, added_at) VALUES('Horror', now());
INSERT INTO genre(name, added_at) VALUES('Drama', now());

--rented_film
CREATE TABLE rented_film (
  id serial primary key,
  user_id INTEGER REFERENCES person(id) ON UPDATE CASCADE ON DELETE CASCADE,
  film_id INTEGER REFERENCES film(id) ON UPDATE CASCADE ON DELETE CASCADE,
  finished INTEGER DEFAULT 0,
  added_at TIMESTAMP
);
