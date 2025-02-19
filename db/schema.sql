-- Warzones
CREATE TABLE IF NOT EXISTS "warzone"(
  "id"      TEXT NOT NULL UNIQUE,
  "server"  INTEGER NOT NULL,
  PRIMARY KEY("id")
);
-- Alliances
CREATE TABLE IF NOT EXISTS "alliance"(
  "id"	        TEXT NOT NULL UNIQUE,
  "warzone_id"  TEXT,
  "tag"	        TEXT NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("warzone_id") REFERENCES "warzone"("id")
);
CREATE TABLE IF NOT EXISTS "alliance_data"(
  "name"	        TEXT,
  "date"              TEXT NOT NULL
    CHECK (
      "date" IS date("date")
    ),
  "power"	        INTEGER,
  "gift_level"	  INTEGER,
  "member_count"	INTEGER,
  "r5_id"	        TEXT,
  "alliance_id"	  TEXT NOT NULL,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("r5_id") REFERENCES "commander"("id")
);
CREATE TABLE IF NOT EXISTS "alliance_alias"(
  "alias"       TEXT,
  "preferred"	  BOOLEAN NOT NULL 
    CHECK (
      "preferred" IN (0, 1)
    ),
  "alliance_id" TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id")
);
-- Commanders
CREATE TABLE IF NOT EXISTS "commander"(
  "id"          TEXT NOT NULL,
  "name"        TEXT,
  "warzone_id"  TEXT,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "commander_data"(
  "date"              TEXT NOT NULL
    CHECK (
      "date" IS date("date")
    ),
  "pfp"	              BLOB,
  "hq_level"	        INTEGER,
  "likes"	            INTEGER,
  "hq_power"	        INTEGER,
  "kills"	            INTEGER,
  "profession_level"	INTEGER,
  "total_hero_power"	INTEGER,
  "alliance_id"	      TEXT,
  "commander_id"	    TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("commander_id") REFERENCES "commander"("id")
);
CREATE TABLE IF NOT EXISTS "commander_alias"(
  "alias"	        TEXT,
  "tag"	          TEXT,
  "preferred"	    BOOLEAN NOT NULL 
    CHECK (
      "preferred" IN (0, 1)
    ),
  "commander_id"	TEXT,
  FOREIGN KEY("commander_id") REFERENCES "commander"("id")
);
-- VsDuels
CREATE TABLE IF NOT EXISTS "vsduel"(
  "id"	          TEXT NOT NULL,
  "date"	        TEXT NOT NULL
    CHECK (
      "date" IS date("date")
    ),
  "league_level"	TEXT NOT NULL,
  "league_id"	    TEXT NOT NULL,
  "tournament_id" TEXT NOT NULL,
  PRIMARY KEY("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_week"(
  "id"            TEXT NOT NULL,
  "vsweek_number"	    TEXT NOT NULL
    CHECK (
      "vsweek_number" IN (0, 1, 2, 3)
    ),
  "vsduel_id"     TEXT NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("vsduel_id") REFERENCES "vsduel"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_alliance"(
  "alliance_id" TEXT,
  "vsduel_week_id" TEXT,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("vsduel_week_id") REFERENCES "vsduel_week"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_day"(
  "name"	      TEXT NOT NULL,
  "short_name"	TEXT NOT NULL,
  "day_of_week"	TEXT NOT NULL
    CHECK (
      "day_of_week" IN (
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday",
        "Saturday",
        "Sunday")
    ),
  "vsduel_points" INTEGER,
  PRIMARY KEY("day_of_week")
);
CREATE TABLE IF NOT EXISTS "vsduel_data"(
  "id"	            TEXT NOT NULL,
  "day_of_week"     TEXT NOT NULL,
  "vsduel_week_id"  TEXT NOT NULL,
  PRIMARY KEY("id"),
  FOREIGN KEY("day_of_week") REFERENCES "vsduel_day"("day_of_week"),
  FOREIGN KEY("vsduel_week_id") REFERENCES "vsduel_week"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_alliance_data"(
  "points"	        INTEGER NOT NULL,
  "vsduel_points"   INTEGER NOT NULL,
  "alliance_id"	    TEXT NOT NULL,
  "vsduel_data_id"	TEXT NOT NULL,
  FOREIGN KEY("alliance_id") REFERENCES "alliance"("id"),
  FOREIGN KEY("vsduel_data_id") REFERENCES "vsduel_data"("id")
);
CREATE TABLE IF NOT EXISTS "vsduel_commander_data"(
  "points"	        INTEGER NOT NULL,
  "rank"	          INTEGER NOT NULL,
  "new"	            BOOLEAN NOT NULL 
    CHECK (
      "new" IN (0, 1)
    ),
  "alliance_id"     TEXT NOT NULL,
  "commander_id"	  TEXT NOT NULL,
  "vsduel_data_id"	TEXT NOT NULL,
  FOREIGN KEY("commander_id") REFERENCES "commander"("id"),
  FOREIGN KEY("vsduel_data_id") REFERENCES "vsduel_data"("id")
);





