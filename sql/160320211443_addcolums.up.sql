alter table Tbl_UserInfo
    modify
        FirstName varchar(100) character set utf8mb4;
alter table Tbl_UserInfo
    modify
        MiddleName varchar(100) character set utf8mb4;
alter table Tbl_UserInfo
    modify
        LastName varchar(100) character set utf8mb4;
alter table Tbl_UserInfo
    modify
        FullName varchar(100) character set utf8mb4;
# alter table accounts.Tbl_UserInfo
#     drop column Gender;
alter table accounts.Tbl_UserInfo
add  column thaiFirstName varchar(100) character set tis620;

alter table accounts.Tbl_UserInfo
add column thaiMiddleName varchar(100) character set tis620;

alter table accounts.Tbl_UserInfo
add column thaiLastName varchar(100) character set tis620;

alter table accounts.Tbl_UserInfo
add column thaiFullName varchar(100) character set tis620;

# alter table accounts.Tbl_UserInfo
#     add column Gender varchar(100) character set tis620;