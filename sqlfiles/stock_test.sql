create table stock_test (
	SID int auto_increment primary key,
	location varchar(20) not null default "",
	model varchar(30) not null default "",
	unit int not null default 0,
	cartons int not null default 0,
	boxes int not null default 0,
	total int not null default 0,
	kind  varchar(20) not null default "",
	notes varchar(200) not null default "",
	update_comments varchar(200) not null default "",
	updated_at timestamp default current_timestamp on update current_timestamp
);

insert into stock_test(location, model, unit,cartons,boxes,total,kind,notes,update_comments) select location,model,unit,cartons,boxes,total,kind,notes,update_comments from stock_updated;