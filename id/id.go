package id

/*
递增id生成器, 预期特征
1. id 单调递增
2. 能进行多实例部署，而不影响id递增的特征
3. 具备冗灾能力

还是 为 n位（持久化） + k位（内存中） + 2位实例号

方案1：
	使用raft实现， 只有master 才能取，其他副本节点，需要转发
方案2：
	raft 实现，
1。每秒统计一次所有实例最大offset， 将offset 上传到master
2。master将统计offset将分发到副本
3。副本更新本地offset
 */
