CREATE TABLE picklist_test (
    PID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    PNO varchar(50) NOT NULL,
    model varchar(50) NOT NULL,
    qty int NOT NULL,
    customer varchar(150) NOT NULL DEFAULT '',
    location varchar(50) NOT NULL,
    status varchar(20) NOT NULL DEFAULT 'Pending',
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

insert into picklist_test(PNO, model,qty,customer,location,status,created_at,updated_at) select PNO,model,qty,customer,location,status,created_at,updated_at from picklist;