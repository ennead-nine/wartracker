--Most Points
SELECT ca.alias AS name, ca.tag AS alliance, MAX(vcd.points) AS points, v.date, vw.vsweek_number, vd.day_of_week
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE ca.preferred=1
GROUP BY name
ORDER BY points DESC
LIMIT 10
	
--P4K Most Points
SELECT ca.alias AS name, ca.tag AS alliance, MAX(vcd.points) AS points, v.date, vw.vsweek_number, vd.day_of_week
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE ca.preferred=1
	AND ca.tag="P4K"
GROUP BY name
ORDER BY points DESC
LIMIT 10

--Most Radar Day Points
SELECT ca.alias AS name, MAX(vcd.points) AS points, v.date, vw.vsweek_number, vd.day_of_week
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE vd.day_of_week="Saturday"
	AND ca.preferred=1
	AND ca.tag="P4K"
GROUP BY name
ORDER BY points DESC
LIMIT 10

--P4K Radar Day Points
SELECT ca.alias AS name, MAX(vcd.points), v.date, vw.vsweek_number, vd.day_of_week
FROM vsduel_commander_data AS vcd
INNER JOIN commander_alias AS ca
	ON vcd.commander_id=ca.commander_id
INNER JOIN vsduel_data AS vd
	ON vcd.vsduel_data_id=vd.id
INNER JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
INNER JOIN vsduel AS v
	ON vw.vsduel_id=v.id
WHERE vd.day_of_week="Wednesday"
AND ca.tag="P4K";


