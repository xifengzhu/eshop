CREATE TABLE casbin_rule (
    p_type VARCHAR(100),
    v0 VARCHAR(100),
    v1 VARCHAR(100),
    v2 VARCHAR(100)
)

INSERT INTO casbin_rule VALUES('p', 'user', 'data', 'read');
-- INSERT INTO casbin_rule(p_type, v0, v1) VALUES('g', 'Bob', 'user');
