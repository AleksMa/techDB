package models

type Status struct {
	Forum  int64 `json:"forum"`
	Post   int64 `json:"post"`
	Thread int64 `json:"thread"`
	User   int64 `json:"user"`
}

const InitScript = `CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    nickname CITEXT,
    about    TEXT,
    email    CITEXT UNIQUE,
    fullname TEXT
);



DROP TABLE IF EXISTS forums CASCADE;
CREATE TABLE forums
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    slug     CITEXT      NOT NULL UNIQUE,
    title    TEXT      NOT NULL,
    authorID BIGINT    NOT NULL,
    FOREIGN KEY (authorID) REFERENCES users (ID) ON DELETE CASCADE
);


DROP TABLE IF EXISTS threads CASCADE;
CREATE TABLE threads
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    created  TIMESTAMP,
    forumID  BIGINT    NOT NULL,
    message  TEXT,

    slug     TEXT,
    title    TEXT,

    authorID BIGINT    NOT NULL,
    FOREIGN KEY (authorID) REFERENCES users (ID) ON DELETE CASCADE,
    FOREIGN KEY (forumID) REFERENCES forums (ID) ON DELETE CASCADE
);


DROP TABLE IF EXISTS posts CASCADE;
CREATE TABLE posts
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    created  TIMESTAMP,
    forumID  BIGINT    NOT NULL,
    isEdited BOOLEAN,
    message  TEXT,
    parentID BIGINT    DEFAULT 0,

    authorID BIGINT    NOT NULL,
    threadID BIGINT    NOT NULL,
    FOREIGN KEY (authorID) REFERENCES users (ID) ON DELETE CASCADE,
    FOREIGN KEY (threadID) REFERENCES threads (ID) ON DELETE CASCADE,
    FOREIGN KEY (forumID) REFERENCES forums (ID) ON DELETE CASCADE
    --FOREIGN KEY (parentID) REFERENCES posts (ID) ON DELETE CASCADE
);


DROP TABLE IF EXISTS votes CASCADE;
CREATE TABLE votes
(
    ID       BIGSERIAL NOT NULL PRIMARY KEY,
    voice    BOOLEAN,
    threadID BIGINT    NOT NULL,
    authorID BIGINT    NOT NULL,
    FOREIGN KEY (authorID) REFERENCES users (ID) ON DELETE CASCADE,
    FOREIGN KEY (threadID) REFERENCES threads (ID) ON DELETE CASCADE
);`
