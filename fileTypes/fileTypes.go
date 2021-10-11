package fileTypes

type FieldTypes int

const Text FieldTypes = 0x00000           //单个字段
const Image FieldTypes = 0x00002          //单个图片
const OnlyHtml FieldTypes = 0x00003       //普通html(不包括图片)
const HtmlWithImage FieldTypes = 0x00004  //html包括图片
const MultipleImages FieldTypes = 0x00005 //多图
const Attr FieldTypes = 0x00006           //标签属性选择器
const Fixed FieldTypes = 0x00007          //固定数据，填什么返回什么,选择器就是返回的数据
const Regular FieldTypes = 0x00008        //正则
