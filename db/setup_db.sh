echo "Waiting for containers to be up"
sleep 10
PARAMS="--insecure --host=database"
SQL="/cockroach/cockroach.sh sql $PARAMS"
$SQL -e "CREATE USER endpoints_admin;"
$SQL -e "SET TIMEZONE = 'America/Bogota';"
$SQL -e "GRANT ALL ON DATABASE defaultdb TO endpoints_admin;"
$SQL -e "CREATE TABLE endpoint_table(dominio varchar NOT NULL,info_servers JSONB NOT NULL,grado_ssl varchar(2) NOT NULL,hora_consulta TIMESTAMP NOT NULL,CONSTRAINT pk_endpoint PRIMARY KEY (dominio));"
echo "SQL Script done"