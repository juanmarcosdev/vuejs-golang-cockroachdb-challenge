CREATE TABLE endpoint_table(
    dominio varchar NOT NULL,
    info_servers JSONB NOT NULL,
    grado_ssl varchar(2) NOT NULL,
    hora_consulta TIMESTAMP NOT NULL,
    CONSTRAINT pk_endpoint PRIMARY KEY (dominio)
);