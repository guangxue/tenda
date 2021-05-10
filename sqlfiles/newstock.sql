create table stock_updated (
	SID int auto_increment primary key,
	location varchar(20) not null default "",
	model varchar(30) not null default "",
	unit int not null default 0,
	cartons int not null default 0,
	boxes int not null default 0,
	total int not null default 0,
	kind  varchar(20) not null default "",
	notes varchar(200) not null default "",
	updated_at timestamp default current_timestamp on update current_timestamp);