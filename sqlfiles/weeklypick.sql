select * from picklist
where created_at BETWEEN '2021-04-26' AND date_add('2021-04-26', INTERVAL 7 DAY);

select model, sum(qty) as total
from weeklypicks_0426
group by  model order by total desc;