--
-- TODO: make payload_hash a fixed size field
-- TODO: additional indexes on symbol name for searching
-- TODO: add comments
--
CREATE TABLE IF NOT EXISTS codeintel_scip_index_documents(
    id SERIAL,
    upload_id integer NOT NULL,
    document_path text NOT NULL,
    payload_hash text NOT NULL,
    PRIMARY KEY (upload_id, document_path)
);

CREATE TABLE IF NOT EXISTS codeintel_scip_documents(
    payload_hash text PRIMARY KEY,
    raw_scip_payload bytea NOT NULL
);

CREATE TABLE IF NOT EXISTS codeintel_scip_symbols(
    upload_id integer NOT NULL,
    symbol_name text NOT NULL,
    index_document_id integer,
    definition_ranges bytea,
    reference_ranges bytea,
    implementation_ranges bytea,
    type_definition_ranges bytea,
    PRIMARY KEY (upload_id, symbol_name, index_document_id)
);
