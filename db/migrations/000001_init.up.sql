CREATE TABLE if not exists songs (
   id SERIAL PRIMARY KEY,
   "group" TEXT NOT NULL,
   song TEXT NOT NULL,
   release_date DATE NOT NULL,
   text TEXT NOT NULL,
   link TEXT NOT NULL
);
