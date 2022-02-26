-- url table
CREATE TABLE "url"
(
    "id"         BIGSERIAL   NOT NULL,
    "user_id"    uuid        NOT NULL,
    "url"        VARCHAR     NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT current_timestamp,
    "deleted_at" TIMESTAMPTZ,
    PRIMARY KEY ("id")
);

CREATE INDEX url_user_id_idx ON url (user_id);

CREATE UNIQUE INDEX url_url_idx ON url (url) WHERE deleted_at IS NULL;
