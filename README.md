# Short-Address

### 主要设计思路

1. 设计HTTP Router和Handler
2. HTTP处理流程中加入middleware
3. 利用Go的Interface来实现可扩展的设计
4. 利用Redis的自增长序列生成短地址
