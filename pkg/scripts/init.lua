print("init")


res =  db_query_dict("select * from tbProducts WHERE PId=?", pid)
if table.getn(res) > 0 then
	product = res[1]
end

table_print(product)

print("DBHost:"..Product.DBHost)

print("State log:")
