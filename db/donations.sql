SELECT DISTINCT ca.alias, d.amount 
FROM donation AS d
JOIN commander_alias AS ca
	ON d.commander_id=ca.commander_id
	AND ca.preferred=1
ORDER BY amount DESC
