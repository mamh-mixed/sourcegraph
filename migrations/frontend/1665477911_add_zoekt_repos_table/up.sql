CREATE TABLE IF NOT EXISTS zoekt_repos (
    repo_id integer NOT NULL PRIMARY KEY,
    commit text,

    CONSTRAINT repo_id_commit_unique UNIQUE (repo_id, commit),

    index_status text DEFAULT 'not_indexed'::text NOT NULL,

    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

INSERT INTO zoekt_repos(repo_id)
SELECT id
FROM repo
LEFT JOIN zoekt_repos zr ON repo.id = zr.repo_id
WHERE zr.repo_id IS NULL
ON CONFLICT (repo_id) DO NOTHING;

CREATE OR REPLACE FUNCTION func_insert_zoekt_repo() RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO zoekt_repos (repo_id) VALUES (NEW.id);

  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trig_zoekt_repos_insert on repo;

CREATE TRIGGER trig_zoekt_repos_insert
AFTER INSERT
ON repo
FOR EACH ROW
EXECUTE FUNCTION func_insert_zoekt_repo();
