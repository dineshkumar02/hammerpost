dbset db pg
dbset bm TPC-C
diset connection pg_sslmode disable
diset connection pg_host {{db_host}}
diset connection pg_port {{db_port}}
diset tpcc pg_superuser {{db_user}}
diset tpcc pg_superuserpass {{db_password}}
diset tpcc pg_storedprocs true
diset tpcc pg_total_iterations {{total_transactions}}
diset tpcc pg_driver timed
diset tpcc pg_rampup {{rampup_duration}}
diset tpcc pg_duration {{test_duration}}
diset tpcc pg_allwarehouse {{all_warehouses}}
diset tpcc pg_count_ware {{warehouses}}
vuset logtotemp 1
loadscript
tcstart
vuset vu {{users}}
vuset delay 0
vucreate
vurun
vudestroy
tcstop