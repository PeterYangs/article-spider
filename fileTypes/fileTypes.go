package fileTypes

type FieldTypes int

const SingleField FieldTypes = 0x00000 //单个字段
const ListField FieldTypes = 0x00001
const SingleImage = 0x00002   //单个图片
const OnlyHtml = 0x00003      //普通html(不包括图片)
const HtmlWithImage = 0x00004 //html包括图片
const ListImages = 0x00005    //多图
