CREATE TABLE
  todos (
    id varchar(128) NOT NULL,
    content varchar(256) NOT NULL,
    done tinyint (1) NOT NULL DEFAULT 0,
    category_id varchar(128) DEFAULT NULL,
    `created_at` int (10) unsigned NOT NULL,
    PRIMARY KEY (id)
  );

CREATE TABLE
  categories (
    id varchar(128) NOT NULL,
    name varchar(256) NOT NULL,
    `created_at` int (10) unsigned NOT NULL,
    PRIMARY KEY (id)
  );