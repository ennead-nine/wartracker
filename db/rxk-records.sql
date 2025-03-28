-- RxK Top 10 All Time
SELECT ca.alias AS cname, MAX(vcd.points) AS points, date(v.date, "+"||((vw.vsweek_number*7)+vdy.day_number)||" days") AS date, vdy.short_name
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE ca.preferred=1
	AND ca.tag="RxK"
GROUP BY cname
ORDER BY points DESC
LIMIT 10

-- RxK Top 10 Daily
SELECT ca.alias AS cname, MAX(vcd.points) AS points, date(v.date, "+"||((vw.vsweek_number*7)+vdy.day_number)||" days") AS date
FROM vsduel_commander_data AS vcd
JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
JOIN vsduel_day AS vdy
	ON vd.day_of_week=vdy.day_of_week
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE ca.preferred=1
	AND ca.tag="RxK"
	AND vdy.short_name="Units"
GROUP BY cname
ORDER BY points DESC
LIMIT 10