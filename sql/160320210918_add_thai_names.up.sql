alter table Tbl_UserInfo
    modify
        FirstName varchar(100) character set tis620
        collate tis620_thai_ci;

alter table Tbl_UserInfo
    modify
        MiddleName varchar(100) character set tis620
        collate tis620_thai_ci;
alter table Tbl_UserInfo
    modify
        LastName varchar(100) character set tis620
        collate tis620_thai_ci;
alter table Tbl_UserInfo
    modify
        FullName varchar(100) character set tis620
        collate tis620_thai_ci;

insert into Tbl_UserInfo
    (UserUUID, IdXref, FirstName, MiddleName, LastName, FullName, Email, UserName, Password)
    value ('hhhh-trhyth-hythyt','idxrefxx687','นี่เป็นประโยคง่ายๆที่คุณต้องแปลงข้อความนี้เป็นภาษาไท','นี่เป็นประโยคง่ายๆที่คุณต้องแปลงข้อความนี้เป็นภาษาไท','นี่เป็นประโยคง่ายๆที่คุณต้องแปลงข้อความนี้เป็นภาษาไท ','นี่เป็นประโยคง่ายๆที่คุณต้องแปลงข้อความนี้เป็นภาษาไท','email@email.com','username@user','this is password')