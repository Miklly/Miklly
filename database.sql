--销售渠道
CREATE TABLE 'ChannelInfo'(
    'ID' varchar(20) primary key,			--微信号
    'Name' varchar(20)						--销售号名称
);
--图片信息
CREATE TABLE 'ImageInfo'(
    'ID' integer primary key AUTOINCREMENT,
    'Length' integer,						--图片文件大小
    'MD5' char(32),							--图片文件MD5
    'SHA1' char(32),						--图片文件SHA1摘要
    'FilePath' varchar(200),				--图片文件路径
    'ThumbnailPath' varchar(200)			--缩略图路径
);
--订单信息
CREATE TABLE 'OrderInfo'(
    'ID' integer primary key AUTOINCREMENT,
    'Name' varchar(20),						--收货人
    'Phone' varchar(15),					--联系电话
    'Address' varchar(1000),				--收货地址
    'ChannelID' varchar(20),				--销售渠道
    'ExpressCompany' varchar(10),			--快递公司
    'ExpressNumber' varchar(20),			--快递单号
    'SendTime' date null,					--发货时间
);
--订单项
CREATE TABLE 'OrderItem'(
    'ID' integer primary key AUTOINCREMENT,
    'OrderID' integer,						--订单编号
    'ImageID' integer,						--图片编号
    'Size' varchar(20),						--尺码
    'SupplierID' integer null,				--供应商编号
    'CreateTime' date,						--下单时间
    'GetTime' date null,					--拿货时间
    'IsSend' bool							--是否发货
);
--供应商
CREATE TABLE 'SupplierInfo'(
    'ID' integer primary key AUTOINCREMENT,
    'Name' varchar(20),						--供应商名称
    'Description' varchar(2000)				--备注
);
--供应商微信平台号
CREATE TABLE 'SupplierWX'(
	'ID' integer primary key AUTOINCREMENT,
	'SupplierID' integer,					--供应商编号
	'WXID' varchar(20),						--微信号
	'Description' varchar(2000)				--备注
)
--供应商朋友圈消息
CREATE TABLE 'SupplierRecord'(
	'ID' integer primary key AUTOINCREMENT,
	'Time' date,							--消息发送时间
	'Title' varchar(1000) null,				--文字标题
	'PicIDs' varchar(100) null,				--图片或视频编号:以逗号(,)分隔
	'CreateTime' date						--消息录入时间
);