-- TESTE 7: Filtrar apenas Entregue
SELECT id, 
       xpath('//ID/text()', xml_documento) as id_xml,
       xpath('//Valor/text()', xml_documento) as valor
FROM relatorio 
WHERE xpath_exists('//Venda[Status="Entregue"]', xml_documento);

-- TESTE 8: Contar por status
SELECT COUNT(*) as total_entregue 
FROM relatorio 
WHERE xpath_exists('//Venda[Status="Entregue"]', xml_documento);

-- TESTE 9: Soma de valores para Entregue
SELECT SUM(CAST((xpath('//Valor/text()', xml_documento))[1]::text AS INTEGER)) as total_valor
FROM relatorio 
WHERE xpath_exists('//Venda[Status="Entregue"]', xml_documento);

-- TESTE 10: Agrupar por status
SELECT 
  (xpath('//Status/text()', xml_documento))[1]::text as status,
  COUNT(*) as quantidade,
  SUM(CAST((xpath('//Valor/text()', xml_documento))[1]::text AS INTEGER)) as total_valor
FROM relatorio
GROUP BY (xpath('//Status/text()', xml_documento))[1]::text;
