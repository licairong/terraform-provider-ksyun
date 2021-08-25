resource "ksyun_redis_instance" "default" {
  available_zone = "cn-beijing-6b"
  name = "redis_2107051609421"
  mode = 2
  capacity = 1
  net_type = 2
  security_group_id = "b6cd4072-4ee0-4fb3-8bd1-2dde798e6bd2,da521976-1e16-48fe-bb3f-7a5d40fb8501"
  vnet_id = "1988b97b-9290-480b-b560-c88a4a90c863"
  vpc_id = "a8979fe2-cf1a-47b9-80f6-57445227c541"
  bill_type = 5
  duration = ""
  pass_word = "Shiwo1101"
  iam_project_id = "0"
  slave_num = 0
  protocol = "4.0"
  reset_all_parameters = false
  parameters = {
    "hash-max-ziplist-entries" ="256",
    "lazyfree-lazy-expire" = "no",
    "lazyfree-lazy-server-del" = "no",
    "slowlog-log-slower-than" = "2000000",
    "slowlog-max-len" = "128"
  }
}