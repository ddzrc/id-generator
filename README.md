本项目包括两种id生成器， 防预测id，自增id，

一. 防预测id
具有以下特征
1.持久化可定制，目前已经实现jdbc
2.防预测 
3.高并发
4.支持多业务隔离


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
                                              
                                              
                                              