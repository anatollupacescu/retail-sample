SELECT
  r.id,
  r.name,
  i.id,
  ri.quantity
FROM
  recipe_ingredient ri,
  recipe r,
  inventory i
WHERE
  ri.recipeid = r.id
  AND ri.inventoryid = i.id;

SELECT * FROM inventory;