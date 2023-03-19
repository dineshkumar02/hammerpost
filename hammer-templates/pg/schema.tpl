dbset db pg
dbset bm TPC-C
diset connection pg_sslmode disable
diset connection pg_host {{db_host}}
diset connection pg_port {{db_port}}
diset tpcc pg_superuser {{db_user}}
diset tpcc pg_superuserpass {{db_password}}

deleteschema
vudestroy

vuset logtotemp 1
diset tpcc pg_storedprocs true
diset tpcc pg_count_ware {{warehouses}}
diset tpcc pg_num_vu {{users}}
buildschema
waittocomplete
quit
