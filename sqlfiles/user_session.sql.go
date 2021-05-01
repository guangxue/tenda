create table user (
	uid int primary key auto_increment,
	email varchar(30) not null default "",
	username varchar(50) not null unique,
	password text not null,
	disabled bool default 0,
	last_login timestamp default current_timestamp);

create table sessions (
	cookie varchar(36) primary key,
	uid int not null,
	session_start_time timestamp not null,
	last_acces timestamp not null,
)