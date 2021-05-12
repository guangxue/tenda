CREATE TABLE last_updated_test (
	LID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
	location varchar(20) NOT NULL,
	model varchar(30) NOT NULL,
	unit int NOT NULL,
	cartons int NOT NULL,
	boxes int NOT NULL,
	total int NOT NULL,
	completed_at timestamp NOT NULL
);

insert into last_updated_test(location,model,unit,cartons,boxes,total,completed_at) select location,model,unit,cartons,boxes,total,completed_at from last_updated;