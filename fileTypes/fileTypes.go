package fileTypes

type FieldTypes int

const SingleField FieldTypes = 1 //单个字段
const ListField FieldTypes = 2
const SingleImage = 3   //单个图片
const OnlyHtml = 4      //普通html(不包括图片)
const HtmlWithImage = 5 //html包括图片
const ListImages = 6    //多图
