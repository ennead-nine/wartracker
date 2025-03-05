--Alliance Data
SELECT 
	vd.day_of_week, 
	a.tag, 
	vad.points, 
	vad.vsduel_points
FROM vsduel_alliance_data AS vad
JOIN alliance AS a
	ON vad.alliance_id=a.id
JOIN vsduel_data AS vd
	ON vad.vsduel_data_id=vd.id
JOIN vsduel_week AS vw
	ON vd.vsduel_week_id=vw.id
WHERE
	vw.vsweek_number=3
	
	