select * from picklist
where created_at BETWEEN '2021-04-26' AND date_add('2021-04-26', INTERVAL 7 DAY);

select model, sum(qty) as total
from weeklypicks_0426
group by  model order by total desc;

SELECT pno, customer,model, qty, updated_at
FROM picklist
WHERE updated_at BETWEEN '2021-08-29%' AND date_add('2021-08-29%', INTERVAL 7 DAY);

SELECT * FROM picklist
WHERE created_at between '2021-10-03' AND date_add('2021-10-03', interval 7 day)
AND model='MW3-3PK' and location='0-G-4';




SELECT PNO, sales_mgr, model, qty, customer, location FROM picklist Where created_at BETWEEN '2021-10-17' AND date_add('2021-10-17', interval 7 day);
