CREATE TABLE IF NOT EXISTS "alliance_data"(
  "date"	TEXT,
  "power"	INTEGER,
  "gift-level"	INTEGER,
  "member-count"	INTEGER,
  "r5-id"	TEXT,
  "alliance-id"	TEXT,
  FOREIGN KEY("alliance-id") REFERENCES "alliance"("id"),
  FOREIGN KEY("r5-id") REFERENCES "commander"("id")
);
CREATE TABLE IF NOT EXISTS "versus-duel"(
  "id"	TEXT NOT NULL,
  "date"	TEXT NOT NULL,
  "alliance1"	TEXT NOT NULL,
  "alliance2"	TEXT NOT NULL,
  "league"	TEXT,
  "week"	INTEGER NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("alliance1") REFERENCES "alliance"("id"),
  FOREIGN KEY("alliance2") REFERENCES "alliance"("id")
);
CREATE TABLE IF NOT EXISTS "alliance"(
  "id"	TEXT NOT NULL,
  "name"	TEXT NOT NULL,
  "tag"	TEXT NOT NULL,
  "server"	INTEGER NOT NULL,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "versus-duel_data"(
  "alliance1-points"	INTEGER,
  "alliance2-points"	INTEGER,
  "versus-day-id"	TEXT,
  "versus-duel-id"	TEXT,
  FOREIGN KEY("versus-day-id") REFERENCES "versus-days"("id"),
  FOREIGN KEY("versus-duel-id") REFERENCES "versus-duel"("id")
);
CREATE TABLE IF NOT EXISTS "versus-duel_commander-data"(
  "points"	INTEGER,
  "rank"	INTEGER,
  "versus-duel-id"	TEXT,
  "commander-id"	TEXT,
  "versus-day-id"	TEXT,
  FOREIGN KEY("commander-id") REFERENCES "commander"("id"),
  FOREIGN KEY("versus-day-id") REFERENCES "versus-days"("id"),
  FOREIGN KEY("versus-duel-id") REFERENCES "versus-duel"("id")
);
CREATE TABLE IF NOT EXISTS "versus-days"(
  "id"	TEXT NOT NULL,
  "name"	TEXT NOT NULL,
  "short-name"	TEXT NOT NULL,
  "day-of-week"	TEXT NOT NULL,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "commander"(
  "id"	TEXT NOT NULL,
  "note-name"	TEXT,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "commander_data"(
  "date"	TEXT NOT NULL,
  "name"	TEXT NOT NULL,
  "kills"	INTEGER,
  "hq-power"	INTEGER,
  "total-hero-power"	INTEGER,
  "commander-id"	TEXT NOT NULL,
  "alliance-id"	TEXT NOT NULL,
  FOREIGN KEY("alliance-id") REFERENCES "alliance"("id"),
  FOREIGN KEY("commander-id") REFERENCES "commander"("id")
);
