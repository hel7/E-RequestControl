CREATE DATABASE IF NOT EXISTS `request_manager`;

use `request_manager`;


CREATE TABLE IF NOT EXISTS Role (
                                    RoleID INT AUTO_INCREMENT PRIMARY KEY,
                                    RoleName VARCHAR(50) NOT NULL UNIQUE,
                                    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
INSERT INTO Role (RoleID, RoleName) VALUES (1, 'Admin');
INSERT INTO Role (RoleID, RoleName) VALUES (2, 'User');

CREATE TABLE User (
                      UserID INT AUTO_INCREMENT PRIMARY KEY,
                      FirstName VARCHAR(100) NOT NULL,
                      LastName VARCHAR(100) NOT NULL,
                      Email VARCHAR(255) NOT NULL UNIQUE,
                      Username VARCHAR(255) UNIQUE,
                      Password VARCHAR(255),
                      RoleID INT NOT NULL,
                      CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                      UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                      FOREIGN KEY (RoleID) REFERENCES Role(RoleID) ON DELETE CASCADE
);

INSERT INTO User (
    FirstName,
    LastName,
    Email,
    Username,
    Password,
    RoleID
)
VALUES (
           'Request',
           'Manager',
           'admin@requestmanager.com',
           'requestmanager_admin',
           'a4dee544b7cac6452eb3e7dc13ebedb058e5b0a6842aba2980922a92108fce25',
           1
       );

CREATE TABLE TicketStatus (
                              StatusID INT AUTO_INCREMENT PRIMARY KEY,
                              Status VARCHAR(50) NOT NULL UNIQUE,
                              CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE Ticket (
                        TicketID INT AUTO_INCREMENT PRIMARY KEY,
                        Title VARCHAR(255) NOT NULL,
                        Description TEXT,
                        StatusID INT NOT NULL,
                        CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                        AssignedTo INT NULL,
                        UserID INT NOT NULL,
                        FOREIGN KEY (StatusID) REFERENCES TicketStatus(StatusID) ON DELETE CASCADE,
                        FOREIGN KEY (AssignedTo) REFERENCES User(UserID) ON DELETE SET NULL,
                        FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
);

CREATE TABLE Notification (
                              NotificationID INT AUTO_INCREMENT PRIMARY KEY,
                              Message TEXT NOT NULL,
                              UserID INT NOT NULL,
                              CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              FOREIGN KEY (UserID) REFERENCES User(UserID) ON DELETE CASCADE
);
