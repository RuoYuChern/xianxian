create DATABASE taoxxxydb;

create table xx_user(
	id bigint generated by default as identity primary key,
	openId varchar(64) NOT NULL,
	userType varchar(12) NOT NULL,
	nickName varchar(64) NOT NULL,
	avatar varchar(256) NOT NULL,
	createTime timestamp NOT NULL,
	updateTime timestamp NOT NULL
);


create unique index oid_type_uq ON xx_user(openId, userType);
create index usr_upt_idx ON xx_user(updateTime);