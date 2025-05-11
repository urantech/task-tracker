CREATE TABLE IF NOT EXISTS users
(
    id       BIGSERIAL PRIMARY KEY,
    email    VARCHAR UNIQUE NOT NULL,
    password VARCHAR        NOT NULL,
    enabled  BOOLEAN DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS user_authority
(
    id        BIGSERIAL PRIMARY KEY,
    authority VARCHAR NOT NULL,
    user_id   BIGINT  NOT NULL
);

CREATE TABLE IF NOT EXISTS task
(
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT  NOT NULL,
    description VARCHAR NOT NULL,
    done        BOOLEAN DEFAULT FALSE
);

ALTER TABLE user_authority
    ADD CONSTRAINT user_authority_user_id_fk FOREIGN KEY (user_id)
        REFERENCES users ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE task
    ADD CONSTRAINT task_user_id_fk FOREIGN KEY (user_id)
        REFERENCES users ON UPDATE RESTRICT ON DELETE RESTRICT;
