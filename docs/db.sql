create database apprank default character set utf8;

/* ランキング情報テーブル */
create table ranks
(
    id bigint not null unique auto_increment,
    updated datetime not null,
    country varchar(64) not null,
    kind int not null,
    rank int not null,
    app_id bigint not null,
    primary key (id)
) engine=innodb;

/* 著作者情報テーブル */
create table artists
(
    id bigint not null unique,
    name text not null,
    url text not null,
    primary key (id)
) engine=innodb;

/* アプリ情報テーブル */
create table apps
(
    id bigint not null unique,
    name text not null,
    url text not null,
    artwork_url text not null,
    kind text not null,
    copyright text not null,
    artist_id bigint not null,
    release_date datetime not null,
    primary key (id)
) engine=innodb;