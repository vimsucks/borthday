CREATE TABLE IF NOT EXISTS birthday (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  uid INTEGER NOT NULL,
  name VARCHAR(255) NOT NULL,
  lunar_birthday DATE NOT NULL,
  solar_birthday DATE NOT NULL
)