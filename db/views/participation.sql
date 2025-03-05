--Member Stuff

--Member List
SELECT ca.alias AS members
FROM commander_alias AS ca
JOIN commander_data AS cd
	ON cd.commander_id=ca.commander_id
WHERE ca.preferred=1
	AND ca.tag="P4K"
ORDER BY ca.alias

--VS Daily Average
SELECT ca.alias AS members, ROUND(AVG(vcd.points)) AS average
FROM commander_alias AS ca
JOIN vsduel_commander_data AS vcd
	ON vcd.commander_id=ca.commander_id
WHERE ca.preferred=1
	AND ca.tag="P4K"
GROUP BY (vcd.commander_id)
ORDER BY ca.alias

--VS Push Day Average
SELECT ca.alias AS members, ROUND(AVG(vcd.points)) AS push_average
FROM commander_alias AS ca
JOIN commander_data AS cd
	ON cd.commander_id=ca.commander_id
JOIN vsduel_commander_data AS vcd
	ON vcd.commander_id=cd.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
WHERE ca.preferred=1
	AND ca.tag="P4K"
	AND vd.push=1
GROUP BY (vcd.commander_id)
ORDER BY ca.alias

--Absent Days
SELECT ca.alias, (COUNT(coalesce(vcd.points, 0)) - (SELECT COUNT(*) FROM vsduel_data)) * -1 AS absent
FROM commander_alias AS ca
INNER JOIN commander_data AS cd
	ON ca.commander_id=cd.commander_id
INNER JOIN alliance AS a
	ON cd.alliance_id=a.id
INNER JOIN vsduel_commander_data AS vcd
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
WHERE ca.preferred=1
	AND ca.commander_id=vcd.commander_id
	AND ca.tag="P4K"
GROUP BY ca.commander_id
ORDER BY ca.alias

--VS Under 7.2 Million
SELECT ca.alias AS members, (COUNT(coalesce(vcd.points, 0)))AS under72
FROM commander_alias AS ca
INNER JOIN commander_data AS cd
	ON ca.commander_id=cd.commander_id
INNER JOIN alliance AS a
	ON cd.alliance_id=a.id
INNER JOIN vsduel_commander_data AS vcd
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
WHERE ca.preferred=1
	AND ca.commander_id=vcd.commander_id
	AND ca.tag="P4K"
	AND vcd.points < 7200000
GROUP BY ca.commander_id
ORDER BY ca.alias

--VS Push Days Under 7.2 Million
SELECT ca.alias AS members, (COUNT(coalesce(vcd.points, 0)))AS push_under72
FROM commander_alias AS ca
INNER JOIN commander_data AS cd
	ON ca.commander_id=cd.commander_id
INNER JOIN alliance AS a
	ON cd.alliance_id=a.id
INNER JOIN vsduel_commander_data AS vcd
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
WHERE ca.preferred=1
	AND ca.commander_id=vcd.commander_id
	AND ca.tag="P4K"
	AND vcd.points < 7200000
	AND vd.push=1
GROUP BY ca.commander_id
ORDER BY ca.alias

--VS Save Day Average
SELECT ca.alias AS members, ROUND(AVG(vcd.points)) AS save_average
FROM commander_alias AS ca
JOIN commander_data AS cd
	ON cd.commander_id=ca.commander_id
JOIN alliance AS a
	ON cd.alliance_id=a.id
JOIN vsduel_commander_data AS vcd
	ON vcd.commander_id=cd.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
WHERE ca.preferred=1
	AND a.tag="P4K"
	AND vd.save=1
GROUP BY (vcd.commander_id)
ORDER BY ca.alias

--Max points
SELECT ca.alias AS members, MAX(vcd.points) AS vsmax,  date(v.date, "+"||((vw.vsweek_number*7)+vdy.day_number)||" days") AS vs_date, vdy.short_name
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON ca.commander_id=vcd.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE ca.tag="P4K"
	AND ca.preferred=1
GROUP BY ca.alias
ORDER BY ca.alias

--Min points
SELECT ca.alias AS members, MIN(vcd.points) AS vsmin,  date(v.date, "+"||((vw.vsweek_number*7)+vdy.day_number)||" days") AS vs_date, vdy.short_name
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON ca.commander_id=vcd.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE ca.tag="P4K"
	AND ca.preferred=1
GROUP BY ca.alias
ORDER BY ca.alias

--Individual VS Day Average
SELECT ca.alias AS members, ROUND(AVG(vcd.points)) AS kill_average
FROM commander_alias AS ca
JOIN vsduel_commander_data AS vcd
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
WHERE ca.preferred=1
	AND ca.tag="P4K"
	AND vdy.short_name="Kill"
GROUP BY (vcd.commander_id)
ORDER BY ca.alias









