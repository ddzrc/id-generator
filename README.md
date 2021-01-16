本项目包括两种id生成器， 防预测id，自增id，

一. 防预测id
具有以下特征
1.持久化可定制，目前已经实现jdbc
2.防预测 
3.高并发
4.支持多业务隔离,
5.支持分布式
使用场景，给用户暴露的订单号

实现原理：

n 持久化 + k （内存中）
         |
         
         
如： 
  持久层存储：1000000 10000
  内存： 1000000 + [0000， 0001, 0002 ..... 9999]， 再次打乱内存中的数据，  1000000 + [8632， 0002, 7875 ..... 9942]
                                                  将 内存中的后缀id放入channnel
                                                  
                                                 ｜
                                               第一次取： 100000 8632
                                               
                                               第二次取： 100000 0002
                                               
                                                 。
                                                 。
                                                 。
                                                 。
                                              超过配置刷新时间， 进行 fresh  -》1000001 + [7442， 0232, 7443 ..... 4775]
                                              
 
                                              超过配置刷新时间， 进行 fresh  -》1000001 + [7442， 0232, 7443 ..... 4775]
                                              
                                              
 二. 分布式自增id（开发中）

         由于分布式自增非常难做， 高并发和严格递增不可兼得，即使在应用场景外得到的是递增id，但是最终落入数据库中的id也不是严格递增的，
 本文设计理念，是使用客户端和服务端共同来维护自增， 客户端有全局变量保存已经获取的最大id， 会轮询实例，全局变量作为参数，进行请求，id生成器实例，会取一个大于参数的值，
 影响自增的一个场景，客户端和服务端同时存在很长一段时间没进行请求或接受请求，相比其他 服务器
 
 
                   实例1                         实例2                                实例master                      实例3                             实例4
                   序列号 -1                   序列号 -2                             协调工作                          序列号 -3                          序列号 - 4
                   
                   
                   


客户端轮询 四个实例，

master：不会直接放回id，而是会进行协调工作
1. 下发，持久化前缀
2. 维护实例个数
3. 调整实例个数
4. 下发实例序号，


实例：
1. 实例的自增值是实例个数，比如 当前id = 1，实例个数是4， 下个id为5
2. 实例序号为，服务启动时向， master申请， 实例对应的id = 序列号 + 实例个数 * n
3. 客户有请求时， 会携带上一个期望值id号， 实例会取下一个大于期望的id,这是维持 id自增的核心


 
 
 
 
 
 
 
 
 
                                              
                                              
