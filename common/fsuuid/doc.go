package fsuuid

/*
UUID(Universally Unique IDentifier)是一个128位数字的唯一标识。RFC4122 描述了具体的规范实现。

1、格式：
	UUID 使用 16 进制表示，共有 36 个字符(32个字母数字+4个连接符"-")，格式为 8-4-4-4-12，如：
		6d25a684-9558-11e9-aa94-efccd7a0659b

2、version1
	[time-low]-[time-mide]-[time-high-and-version]-[clock-seq-and-reserved 和 clock-seq-low]-[node]
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                          time_low                             |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|       time_mid                |         time_hi_and_version   |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|clk_seq_hi_res |  clk_seq_low  |         node (0-1)            |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
		|                         node (2-5)                            |
		+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

	time-low = 32位 unsigned integer (时间戳的低32位)
	time-mid = 16位 unsigned integer (时间戳的中间 16 位)
	time-high-and-version = 16位 unsigned integer (时间戳高位部分与版本(Version)号混合)

	clock-seq-and-reserved = 8位 unsigned integer (时钟序列高位部分与预定义变量(Variant)混合组成)
	clock-seq-low = 8位 unsigned integer

	node = 48位 unsigned integer

	如：
		xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx

		M 中使用 4 位来表示 UUID 的版本。上面第三部分，先用时间戳的高位填满，然后再把 M 处修改为版本号，
		因此 v1 的 M 位总是 0x1

		N 中使用 1-3 位表示不同的 variant，这里采用 RFC4122 标准，因此 N 位的二进制表示总是：10XX

		最后 node 部分，规范推荐使用网卡 MAC 地址，这里直接使用随机数表示

	版本号(version), 4 bits, 一共以下5个版本:
		0001 时间的版本
		0010 DCE Security
		0011 MD5哈希
		0100 (伪)随机数
		0101 SHA-1哈希

	变量(variant), 或称做类型, 4 bits, 包括以下4种（其中X为任意值）:
		0XX NCS兼容预留
		10X RFC4122采用
		110 微软兼容预留
		111 还未定义, 留作以后它用

3、version2
	暂时不实现

4、version3
	暂时不实现

5、version4
	xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx

	该版本下:
		M 位(version 位)总是 0x4（即二进制的：0100），表示（伪）随机数
		N 位总是二进制的：10XX

6、version5
	暂时不实现

*/
