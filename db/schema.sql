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
CREATE TABLE IF NOT EXISTS "vsduel_data"(
  "alliance1-points"	INTEGER,
  "alliance2-points"	INTEGER,
  "vsduel_day-id"	TEXT,
  "vsduel-id"	TEXT,
  FOREIGN KEY("vsduel-id") REFERENCES "vsduel"("id"),
  FOREIGN KEY("vsduel_day-id") REFERENCES "vsduel_day"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_day"(
  "id"	TEXT NOT NULL,
  "name"	TEXT NOT NULL,
  "short-name"	TEXT NOT NULL,
  "day-of-week"	TEXT NOT NULL,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_commanders"(
  "points"	INTEGER,
  "rank"	INTEGER,
  "vsduel-id"	TEXT,
  "commander-id"	TEXT,
  "vsduel_day-id"	TEXT,
  FOREIGN KEY("commander-id") REFERENCES "commander"("id"),
  FOREIGN KEY("vsduel-id") REFERENCES "vsduel"("id"),
  FOREIGN KEY("vsduel_day-id") REFERENCES "vsduel_day"("id")
);
CREATE TABLE IF NOT EXISTS "alliance"(
  "id"	TEXT NOT NULL,
  "server"	INTEGER NOT NULL,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "vsduel"(
  "id"	TEXT NOT NULL,
  "date"	TEXT NOT NULL,
  "alliance1-id"	TEXT NOT NULL,
  "alliance2-id"	TEXT NOT NULL,
  "league"	TEXT,
  "week"	INTEGER NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("alliance1-id") REFERENCES "alliance"("id"),
  FOREIGN KEY("alliance2-id") REFERENCES "alliance"("id")
);
CREATE TABLE IF NOT EXISTS "alliance_data"(
  "name"	TEXT,
  "tag"	TEXT,
  "date"	TEXT NOT NULL UNIQUE,
  "power"	INTEGER,
  "gift_level"	INTEGER,
  "member_count"	INTEGER,
  "r5_id"	TEXT,
  "alliance_id"	TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("r5_id") REFERENCES "commander"("id")
);
