CREATE TABLE IF NOT EXISTS "commander"(
  "id"	TEXT NOT NULL,
  "note_name"	TEXT,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "commander_data"(
  "date"	TEXT NOT NULL,
  "pfp"	BLOB,
  "hq_level"	INTEGER,
  "likes"	INTEGER,
  "hq_power"	INTEGER,
  "kills"	INTEGER,
  "profession_level"	INTEGER,
  "total_hero_power"	INTEGER,
  "commander_id"	TEXT,
  "alliance_id"	TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("commander_id") REFERENCES "commander"("id")
);
CREATE TABLE IF NOT EXISTS "alliance_data"(
  "name"	TEXT,
  "date"	TEXT NOT NULL,
  "power"	INTEGER,
  "gift_level"	INTEGER,
  "member_count"	INTEGER,
  "r5_id"	TEXT,
  "alliance_id"	TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("r5_id") REFERENCES "commander"("id")
);
CREATE TABLE IF NOT EXISTS "alliance"(
  "id"	TEXT NOT NULL UNIQUE,
  "server"	INTEGER NOT NULL,
  "tag"	TEXT NOT NULL,
  PRIMARY KEY("id","server","tag")
);
CREATE TABLE IF NOT EXISTS "vsduel"(
  "id"	TEXT NOT NULL,
  "date"	TEXT NOT NULL,
  "league"	TEXT,
  "week"	INTEGER,
  PRIMARY KEY("id","date")
);
CREATE TABLE IF NOT EXISTS "vsduel_day"(
  "id"	TEXT NOT NULL,
  "name"	TEXT NOT NULL,
  "short_name"	TEXT NOT NULL,
  "day_of_week"	TEXT NOT NULL,
  PRIMARY KEY("id","name")
);
CREATE TABLE IF NOT EXISTS "vsduel_data"(
  "id"	TEXT NOT NULL,
  "vsduel_day_id"	TEXT NOT NULL,
  "vsduel_id"	TEXT NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("vsduel_day_id") REFERENCES "vsduel_day"("id"),
  FOREIGN KEY("vsduel_id") REFERENCES "vsduel"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_alliance"(
  "points"	INTEGER,
  "tag"	TEXT,
  "alliance_id"	TEXT,
  "vsduel_data_id"	TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("vsduel_data_id") REFERENCES "vsduel_data"
);
CREATE TABLE IF NOT EXISTS "vsduel_commander"(
  "points"	INTEGER NOT NULL,
  "rank"	INTEGER NOT NULL,
  "name"	TEXT,
  "commander_id"	TEXT NOT NULL,
  "vsduel_data_id"	TEXT NOT NULL,
  FOREIGN KEY("commander_id") REFERENCES "commander"("id"),
  FOREIGN KEY("vsduel_data_id") REFERENCES ""
);
