CREATE TABLE lastest_stock (
	SID INT AUTO_INCREMENT PRIMARY KEY,
	location VARCHAR(20),
	model VARCHAR(50),
	unit INT NOT NULL,
	cartons INT NOT NULL,
	boxes INT DEFAULT 0,
	total INT NOT NULL,
	description TEXT,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);