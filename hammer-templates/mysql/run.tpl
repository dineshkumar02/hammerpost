dbset db mysql
dbset bm TPC-C
diset connection mysql_host {{db_host}}
diset connection mysql_port {{db_port}}
diset connection mysql_socket null
diset tpcc mysql_user {{db_user}}
diset tpcc mysql_pass {{db_password}}
diset tpcc mysql_count_ware {{warehouses}}
diset tpcc mysql_partition false
diset tpcc mysql_storage_engine innodb


diset tpcc mysql_total_iterations {{total_transactions}}
diset tpcc mysql_driver timed
diset tpcc mysql_rampup {{rampup_duration}}
diset tpcc mysql_duration {{test_duration}}
diset tpcc mysql_allwarehouse {{all_warehouses}}
vuset logtotemp 1
loadscript
tcstart
vuset vu {{users}}
vuset delay 0
vucreate
vurun
vudestroy
tcstop