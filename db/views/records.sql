--Most Points
SELECT ca.alias AS cname, ca.tag AS alliance, MAX(vcd.points) AS points, date(v.date, "+"||((vw.vsweek_number*7)+vdy.day_number)||" days") AS vs_date, vdy.short_name
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
INNER JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE ca.preferred=1
	AND v.league_id LIKE "11-%"
GROUP BY cname
ORDER BY points DESC
LIMIT 10
	
--P4K Most Points
SELECT ca.alias AS cname, MAX(vcd.points) AS points, date(v.date, "+"||((vw.vsweek_number*7)+vdy.day_number)||" days") AS vs_date, vdy.short_name
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
INNER JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE ca.preferred=1
	AND ca.tag="P4K"
GROUP BY cname
ORDER BY points DESC
LIMIT 10

--Most  Day Points
SELECT ca.alias AS cname, ca.tag AS alliance, MAX(vcd.points) AS points, date(v.date, "+"||((vw.vsweek_number*7)+vdy.day_number)||" days") AS vs_date, vdy.short_name
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
INNER JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE vd.day_of_week="Wednesday"
	AND ca.preferred=1
	AND v.league_id LIKE "11-%"
GROUP BY cname
ORDER BY points DESC
LIMIT 10

--P4K  Day Points
SELECT ca.alias AS cname, MAX(vcd.points) AS points, date(v.date, "+"||((vw.vsweek_number*7)+vdy.day_number)||" days") AS vs_date, vdy.short_name
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
INNER JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE vd.day_of_week="Saturday"
	AND ca.preferred=1
	AND ca.tag="P4K"
GROUP BY cname
ORDER BY points DESC
LIMIT 10


