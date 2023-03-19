dbset db mysql
dbset bm TPC-C
diset connection mysql_host {{db_host}}
diset connection mysql_port {{db_port}}
diset connection mysql_socket null
diset tpcc mysql_user {{db_user}}
diset tpcc mysql_pass {{db_password}}

deleteschema
vudestroy

vuset logtotemp 1
diset tpcc mysql_count_ware {{warehouses}}
diset tpcc mysql_partition true
diset tpcc mysql_num_vu {{users}}
diset tpcc mysql_storage_engine innodb
buildschema
waittocomplete
quit