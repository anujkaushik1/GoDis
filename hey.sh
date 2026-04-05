for i in {1..20000}; do
  redis-cli -p 7379 SET "key_$i" "value_$i"
done