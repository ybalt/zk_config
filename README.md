# zk_config
small go program that reads Zookeper node with key=value lines and replace key in input file to new output file


run:
zk_config -zk <zookeeper host:port> -regex <brackets> -path <zk node path> -in <input file> -o <output file>

zk node should be in form

var1=val1
var2=val2

i.e. key=value pairs with newline as delimiter of lines

input file may have any templates like 
value: {{var1}}

so it will be replaced in output file
value: val1

regexp for double brackets {{value}}:
-regexp "\{{(.*?)\}}"

for single {value}:
-regexp "\{(.*?)\}"
