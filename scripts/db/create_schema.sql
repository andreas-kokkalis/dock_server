CREATE TYPE enum_admin_status AS ENUM('active', 'deleted');
CREATE TABLE admins(
    id SERIAL PRIMARY KEY,
    username varchar(60) NOT NULL UNIQUE,
    password varchar(100) NOT NULL,
	status enum_admin_status NOT NULL DEFAULT 'active'
);
