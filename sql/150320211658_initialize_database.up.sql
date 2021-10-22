create schema if not exists accounts;
use accounts;
set foreign_key_checks = 0;
drop table if exists UserInfo;
drop table if exists AddressInfo;
drop table if exists ContactInfo;
set foreign_key_checks = 1;
create table UserInfo
(
    tab_id     varchar(100)                        null,
    id         int auto_increment,
    key (id),
    UserUUID   varchar(100) primary key            not null,
    IdXref     varchar(100) unique,
    FirstName  varchar(30),
    MiddleName varchar(30),
    LastName   varchar(30),
    FullName   varchar(100),
    Email      varchar(30)                         not null,
    UserName   varchar(30),
    Password   varchar(500),
    CreatedAt  timestamp default current_timestamp not null,
    UpdatedAt  timestamp on update current_timestamp,
    DeletedAt  timestamp                           null,
    Status     varchar(10)
);
create table AddressInfo
(
    tab_id   varchar(30)  null,
    id       int primary key auto_increment,
    UserUUID varchar(100) not null,
    Address1 varchar(100) not null,
    Address2 varchar(100),
    Address3 varchar(100),
    Area     varchar(100),
    City     varchar(100),
    State    varchar(100) not null,
    Country  varchar(100),
    PinCode  varchar(100),
    foreign key (UserUUID) references UserInfo (UserUUID)
);

create table ContactInfo
(
    tab_id          varchar(100) null,
    id              int primary key auto_increment,
    UserUUID        varchar(100) not null,
    SocialMediaID   varchar(100),
    WebSite         varchar(100),
    DialCode        varchar(100),
    PhoneNumber     varchar(100) not null,
    FullPhoneNumber varchar(100),
    foreign key (UserUUID) references UserInfo (UserUUID)
);
drop table if exists Tbl_uploadhistory;
create table Tbl_uploadhistory
(
    tab_id        varchar(100) null,
    TransactionId int auto_increment primary key,
    FileName      varchar(100),
    StartTime     timestamp,
    EndTime       timestamp,
    TotalTime     varchar(100)
);
# drop table if exists Tbl_Result;
# create table if not exists Tbl_Subject
# (
#     tab_id  varchar(100) null,
#     Chemistry varchar(100),
#
# );


