CREATE EXTENSION IF NOT EXISTS CITEXT;

DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS thread CASCADE;
DROP TABLE IF EXISTS users_forum;
DROP TABLE IF EXISTS post CASCADE;
DROP TABLE IF EXISTS vote CASCADE;

CREATE UNLOGGED TABLE users (
                       nickname CITEXT UNIQUE PRIMARY KEY,
                       fullname TEXT NOT NULL,
                       about TEXT,
                       email CITEXT NOT NULL UNIQUE
);

CREATE UNLOGGED TABLE forum (
                       title TEXT,
                       "user" CITEXT,
                       slug CITEXT PRIMARY KEY UNIQUE,
                       posts BIGINT DEFAULT 0,
                       threads INT DEFAULT 0,
                       FOREIGN KEY ("user") REFERENCES users(nickname)
);


CREATE UNLOGGED TABLE thread (
                        id SERIAL PRIMARY KEY,
                        title TEXT,
                        author CITEXT,
                        forum CITEXT,
                        message TEXT,
                        votes INT DEFAULT 0,
                        slug CITEXT UNIQUE,
                        created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                        FOREIGN KEY (author) REFERENCES "users"(nickname),
                        FOREIGN KEY (forum)  REFERENCES "forum" (slug)
);

CREATE unlogged TABLE users_forum (
                             nickname CITEXT NOT NULL,
                             slug CITEXT NOT NULL,
                             FOREIGN KEY (nickname) REFERENCES users(nickname),
                             FOREIGN KEY (slug) REFERENCES forum (slug),
                             UNIQUE (nickname, slug)
);

CREATE UNLOGGED TABLE post(
                     id BIGSERIAL PRIMARY KEY,
                     parent BIGINT DEFAULT 0,
                     author CITEXT,
                     message TEXT,
                     isEdited BOOLEAN DEFAULT FALSE,
                     forum CITEXT,
                     thread INT,
                     created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                     paths BIGINT[] DEFAULT ARRAY []::INTEGER[],
                     FOREIGN KEY (author) REFERENCES users(nickname),
                     FOREIGN KEY (forum) REFERENCES forum(slug),
                     FOREIGN KEY (thread) REFERENCES thread(id)
);

CREATE  UNLOGGED TABLE vote (
                      id BIGSERIAL PRIMARY KEY,
                      nickname CITEXT,
                      voice INT,
                      thread INT NOT NULL,
                      FOREIGN KEY (nickname) REFERENCES users(nickname),
                      FOREIGN KEY (thread) REFERENCES thread(id),
                      UNIQUE (thread, nickname)
);

CREATE OR REPLACE FUNCTION add_votes() RETURNS TRIGGER AS
$$
BEGIN
    UPDATE thread
    SET votes=(votes + NEW.voice)
    WHERE id = NEW.thread;
    return NEW;
end
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert_vote
    AFTER INSERT
    ON vote
    FOR EACH ROW
    EXECUTE PROCEDURE add_votes();


CREATE OR REPLACE FUNCTION update_thread_votes() RETURNS TRIGGER AS
$$
BEGIN
    IF OLD.voice <> NEW.voice THEN
        UPDATE thread
        SET votes=(votes + NEW.voice * 2)
        WHERE id = NEW.thread;
    END IF;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_update_voice
    AFTER UPDATE
    ON vote
    FOR EACH ROW
    EXECUTE PROCEDURE update_thread_votes();


CREATE OR REPLACE FUNCTION new_user_forum() RETURNS TRIGGER AS 
$$
BEGIN
    INSERT INTO users_forum (nickname, slug)
    VALUES (new.author, new.forum)
    ON CONFLICT DO NOTHING;
    RETURN new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert_thread_update_user
    AFTER INSERT
    ON thread
    FOR EACH ROW
    EXECUTE PROCEDURE new_user_forum();

CREATE TRIGGER after_insert_post
    AFTER INSERT
    ON post
    FOR EACH ROW
    EXECUTE PROCEDURE new_user_forum();


CREATE OR REPLACE FUNCTION update_paths_post() RETURNS TRIGGER AS
$$
DECLARE
    parent_path         BIGINT[];
    first_parent_thread INT;
BEGIN
    IF (NEW.parent = 0) THEN
        NEW.paths := array_append(NEW.paths, NEW.id);
    ELSE
        SELECT paths FROM post WHERE id = NEW.parent INTO parent_path;
        SELECT thread FROM post WHERE id = parent_path[1] INTO first_parent_thread;

        IF NOT FOUND OR first_parent_thread <> NEW.thread THEN
            RAISE EXCEPTION 'parent post was created in another thread'
            USING ERRCODE = '77777';
        END IF;

        NEW.paths := NEW.paths || parent_path || NEW.id;
    END IF;

    UPDATE forum
    SET posts=posts + 1
    WHERE forum.slug = NEW.forum;
    RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_post
    BEFORE INSERT
    ON post
    FOR EACH ROW
    EXECUTE PROCEDURE update_paths_post();


CREATE OR REPLACE FUNCTION increment_counter_threads() RETURNS TRIGGER AS 
$$
BEGIN
    UPDATE forum
    SET threads = forum.threads + 1
    WHERE slug = NEW.forum;
RETURN NEW;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_insert_thread
    AFTER INSERT
    ON thread
    FOR EACH ROW
    EXECUTE PROCEDURE increment_counter_threads();

CREATE INDEX IF NOT EXISTS idx_users_nickname ON users (nickname);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE INDEX IF NOT EXISTS idx_forum_slug ON forum (slug);

CREATE INDEX IF NOT EXISTS idx_thread_slug ON thread (slug);
CREATE INDEX IF NOT EXISTS idx_thread_forum ON thread (forum);
CREATE INDEX IF NOT EXISTS idx_thread_created ON thread (created);

CREATE INDEX IF NOT EXISTS idx_users_forum_nickname_slug ON users_forum(nickname, slug);
CREATE INDEX IF NOT EXISTS idx_users_forum_nickname ON users_forum(nickname);
CREATE INDEX IF NOT EXISTS idx_users_forum_slug ON users_forum(slug);

CREATE INDEX IF NOT EXISTS idx_post_thread_paths_id ON post (thread, paths, id);
CREATE INDEX IF NOT EXISTS idx_post_thread_id_paths1_parent ON post (thread, (paths[1]), parent);
CREATE INDEX IF NOT EXISTS idx_paths1_id on post ((paths[1]), id);
CREATE INDEX IF NOT EXISTS idx_post_paths1_paths_id ON post ((paths[1]), paths, id);
