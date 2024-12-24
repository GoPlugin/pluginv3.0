echo "<<<<<<<<<------------------STARTING PLUGIN NODE--------------------->>>>>>>>>" 
plugin --admin-credentials-file apicredentials.txt -c config.toml -s secrets.toml node start 
echo "<<<<<<<<<<<-------------------Plugin node is running .. use "pm2 status" to check the status-------->>>>>>>>>>>"

