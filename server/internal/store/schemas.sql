CREATE TABLE IF NOT EXISTS [donations] (
    [id] INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    [date] DATETIME,
    [name] TEXT,
    [type] TEXT,
    [quantity] INTEGER,
    [description] TEXT,
);
CREATE TABLE IF NOT EXISTS [distributions] (
    [id] INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    [donation_id] INTEGER,
    [date] DATETIME
    [type] TEXT,
    [quantity] INTEGER,
    [description] TEXT,
    FOREIGN KEY(donation_id) REFERENCES donations(id)
);