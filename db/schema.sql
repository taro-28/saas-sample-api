CREATE TABLE
  todos (
    id varchar(128) NOT NULL,
    content varchar(256) NOT NULL,
    done tinyint (1) NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
  );