saga 模式还是比较好理解的
将一个大事务拆分为多个小事务
每一个小事务对于一个补偿，大事务执行到哪里出错就从哪开始补偿。