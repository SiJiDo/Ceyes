# Ceyes
一款基于fofa根据域名或fofa语法收集C段分布数量的工具。

### 2023.5.29 更新
1.修复了count计数去重存在问题的bug

2.添加-cloud参数，可识别该段是否为云服务器

云服务器的fofa指纹目前整理如下，主要识别方式通过org和cloud_name 2个字段进行判别，命中cloud_name基本上确定是云了，如果命中org字段可能是云也有可能是云厂商自家的资产。
```
华为云 - Huawei Cloud Service data center [cloud_name="HuaWeiCloud"]

阿里云 - Hangzhou Alibaba Advertising Co.,Ltd. & Alibaba US Technology Co., Ltd. [cloud_name="aliyun"]

腾讯云 - Shenzhen Tencent Computer Systems Company Limited & Tencent Building, Kejizhongyi Avenue [cloud_name="tencent"]

百度云 - Beijing Baidu Netcom Science and Technology Co., Ltd. [cloud_name="baidu"]

亚马逊云 - AMAZON-02 & AMAZON-AES [cloud_name="Amazon"]

google云 - GOOGLE-IT & GOOGLE [cloud_name="google"]

azure云 - MICROSOFT-CORP-MSN-AS-BLOCK [cloud_name="azure"]

Cloudflare - CLOUDFLARENET [cloud_name="Cloudflare"]
```

效果如下

```
main.exe -f 1.txt -sc -cloud
```

![5](img/5.png)


### 编译&运行
```
git clone https://github.com/SiJiDo/Ceyes.git
cd Ceyes/main
go build
```
之后会在目标下产生一个main的可执行文件，第一次使用直接运行它，会在相同目录下生成一个config.yaml文件
内容如下
```
fofa_email: ""
Fofa_api: ""
```
填上自己的fofa邮箱和api即可，这里最好用高级会员，因为高级会员api一次可以查1W条数据，方便更加全面的统计c段分布。如果是普通会员的api查的数据很少，不一定全面。

### 使用方式
功能点如下
```
main.exe -h
Usage of main.exe:
  -cloud
        check cloud
  -d string
        domain deafult use dork like (domain= xxx || host= xxx)
  -f string
        domain text
  -s string
        fofa search dork
  -sc
        sort result by count, deafult ip sort
```
#### 排序功能

通过-sc参数可以对结果进行排序，加上-sc参数会将结果按c段个数排序，不加-sc则按ip进行排序

#### 用法1 根据fofa语句查询分布情况

这里得注意下，因为fofa语法中要用到双引号，所有正确写法如下
```
main.exe -s "title=\"北京大学\"" -sc
```

![1](img/1.png)

#### 用法2 根据主域名查询分布情况

```
main.exe -d "jd.com" -sc
```

![2](img/2.png)

#### 用法3 通过文件进行域名查询

这里文件内容通常是收集到主域名，比如下

```
//1.txt

qq.com
weixin.com
tencent.com
```

Ceyes工具会对目标进行每个域名的ip分布进行信息收集，也会对最后的结果进行汇总

![image-20230529101524961](img/4.png)

当然我之前写到IEyes可以通过查备案快速收集域名，配合起来用也很不错

```
https://github.com/SiJiDo/IEyes
```

![image-20230529101524961](img/3.png)

