CREATE TABLE IF NOT EXISTS "alliance"(
  "id"	TEXT NOT NULL,
  "server"	INTEGER NOT NULL,
  PRIMARY KEY("id")
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
CREATE TABLE IF NOT EXISTS "commander"(
  "id"	TEXT NOT NULL,
  "note_name"	TEXT,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "commander_data"(
  "date"	TEXT NOT NULL,
  "name"	TEXT NOT NULL,
  "hq_power"	INTEGER,
  "kills"	INTEGER,
  "profession_level"	INTEGER,
  "total_hero_power"	INTEGER,
  "commander_id"	TEXT NOT NULL,
  "alliance_id"	TEXT NOT NULL,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("commander_id") REFERENCES "commander"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel"(
  "id"	TEXT NOT NULL,
  "date"	TEXT NOT NULL,
  "league"	TEXT,
  "week"	INTEGER NOT NULL,
  "alliance1_id"	TEXT NOT NULL,
  "alliance2_id"	TEXT NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("alliance1_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("alliance2_id") REFERENCES "alliance"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_commanders"(
  "points"	INTEGER,
  "rank"	INTEGER,
  "vsduel_id"	TEXT,
  "commander_id"	TEXT,
  "vsduel_day_id"	TEXT,
  FOREIGN KEY("commander_id") REFERENCES "commander"("id"),
  FOREIGN KEY("vsduel_day_id") REFERENCES "vsduel_day"("id"),
  FOREIGN KEY("vsduel_id") REFERENCES "vsduel"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_data"(
  "alliance1_points"	INTEGER,
  "alliance2_points"	INTEGER,
  "vsduel_day_id"	TEXT,
  "vsduel_id"	TEXT,
  FOREIGN KEY("vsduel_day_id") REFERENCES "vsduel_day"("id"),
  FOREIGN KEY("vsduel_id") REFERENCES "vsduel"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_day"(
  "id"	TEXT NOT NULL,
  "name"	TEXT NOT NULL,
  "short_name"	TEXT NOT NULL,
  "day_of_week"	TEXT NOT NULL,
  PRIMARY KEY("id")
);
