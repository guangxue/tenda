create table stock_0430 (
	SID int auto_increment primary key,
	location varchar(20) not null default "",
	model varchar(30) not null default "",
	unit int not null default 0,
	cartons int not null default 0,
	boxes int not null default 0,
	total int not null default 0,
	updated_at timestamp default current_timestamp on update current_timestamp);


insert into stock_0430 (location,model,unit,cartons,boxes,total) select location, model, unit, cartons, boxes, total from stock_updated;