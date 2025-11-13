
-- applicant_query.sql
-- SQL Query for "Advertising System Failures Report"


SELECT
    -- Concatenar el nombre y apellido de los clientes
    c.first_name || ' ' || c.last_name AS customer,
    -- Contar cantidad de eventos
    COUNT(*) AS failures
FROM
    -- Tabla de clientes, se asigna alias 
    customers AS c
JOIN
    -- Tabla de campañas, se asigna alias que sirve para unir cada cliente con sus campañas
    campaigns AS ca ON c.id = ca.customer_id
JOIN
    -- Tabla de eventos, se asigna alias que sirve para unir las campañas con los eventos que tuevieron
    events AS e ON ca.id = e.campaign_id
WHERE
    -- Filtra los eventos que fallaron 'failure'
    e.status = 'failure'
GROUP BY
    -- Agrupa los resultados por cada cliente
    c.first_name, c.last_name
HAVING
    -- Se usa HAVING que funciona com un WHERE pero para grupos y solo se muestran los clientes con mas de 3 events fallidos
    COUNT(*) > 3
ORDER BY
    -- Ordena los resultados por la cntidad de fallas de mayor a menos
    failures DESC;
