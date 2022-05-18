<template>
  <div class="tinymce-box">
    <Editor
      v-model="content"
      :init="init"
      :disabled="disabled"
    />
  </div>
</template>

<script>
export default {
  name: 'TEditor',
}
</script>

<script setup>
import { ref } from 'vue'
import tinymce from 'tinymce/tinymce' // tinymce默认hidden，不引入不显示
import Editor from '@tinymce/tinymce-vue'
import 'tinymce/themes/silver/theme' // 主题文件
import 'tinymce/icons/default'
import 'tinymce/models/dom'
// tinymce插件可按自己的需要进行导入
// 更多插件参考：https://www.tiny.cloud/docs/plugins/
import 'tinymce/plugins/image' // 插入上传图片插件
import 'tinymce/plugins/importcss' // 图片工具
import 'tinymce/plugins/media' // 插入视频插件
import 'tinymce/plugins/table' // 插入表格插件
import 'tinymce/plugins/lists' // 列表插件
import 'tinymce/plugins/charmap' // 特殊字符
import 'tinymce/plugins/wordcount' // 字数统计插件
import 'tinymce/plugins/codesample' // 插入代码
import 'tinymce/plugins/code' // 查看源码
import 'tinymce/plugins/fullscreen' // 全屏
import 'tinymce/plugins/link' //
import 'tinymce/plugins/preview' // 预览
import 'tinymce/plugins/template' // 插入模板
import 'tinymce/plugins/save' // 保存
import 'tinymce/plugins/searchreplace' // 查询替换
import 'tinymce/plugins/pagebreak' // 分页
import 'tinymce/plugins/insertdatetime'// 时间插入
import { computed, onMounted, getCurrentInstance } from 'vue'

const context = getCurrentInstance()

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  disabled: {
    type: Boolean,
    default: true,
  },
  plugins: {
    type: [String, Array],
    default:
        'lists image media table wordcount save preview',
  },
  toolbar: {
    type: [String, Array],
    default:
        'formats undo redo paste print fontsizeselect fontselect template fullpage|wordcount ltr rtl visualchars visualblocks toc spellchecker searchreplace|save preview pagebreak nonbreaking|media image|outdent indent aligncenter alignleft alignright alignjustify lineheight  underline quicklink h2 h3 blockquote numlist bullist table removeformat forecolor backcolor bold italic  strikethrough hr charmap link insertdatetime|subscript superscript cut codesample code |anchor preview fullscreen|help',
  },
  height: {
    type: Number,
    default: 500,
  },
  width: {
    type: Number,
    default: 680,
  },
})

const init = ref({
  language_url: '/tinymce/langs/zh-Hans.js', // 引入语言包文件
  language: 'zh_CN', // 语言类型
  width: props.width,
  skin_url: '/tinymce/skins/ui/oxide', // 皮肤：浅色
  // skin_url: '/tinymce/skins/ui/oxide-dark',//皮肤：暗色

  plugins: props.plugins, // 插件配置
  // toolbar: props.toolbar, //工具栏配置，设为false则隐藏
  toolbar_mode: 'sliding',
  menubar: false, // 菜单栏配置，设为false则隐藏，不配置则默认显示全部菜单，也可自定义配置--查看 http://tinymce.ax-z.cn/configure/editor-appearance.php --搜索“自定义菜单”
  menu: {
    // file: { title: '文件', items: 'newdocument' },
    edit: { title: '编辑', items: 'undo redo | cut copy paste pastetext | selectall' },
    insert: { title: '插入', items: 'link image  |  hr' },
    view: { title: '查看', items: 'visualaid' },
    // format: {
    //   title: '格式',
    //   items:
    //     'bold italic underline strikethrough superscript subscript | formats | removeformat',
    // },
    // table: { title: '表格', items: 'inserttable tableprops deletetable | cell row column' },
    // tools: { title: '工具', items: 'spellchecker code' },
  },
  fontsize_formats: '12px 14px 16px 18px 20px 22px 24px 28px 32px 36px 48px 56px 72px', // 字体大小
  font_formats:
      '微软雅黑=Microsoft YaHei,Helvetica Neue,PingFang SC,sans-serif;苹果苹方=PingFang SC,Microsoft YaHei,sans-serif;宋体=simsun,serif;仿宋体=FangSong,serif;黑体=SimHei,sans-serif;Arial=arial,helvetica,sans-serif;Arial Black=arial black,avant garde;Book Antiqua=book antiqua,palatino;',

  height: props.height, // 注：引入autoresize插件时，此属性失效
  placeholder: '在这里输入文字',
  branding: false, // tiny技术支持信息是否显示
  resize: false, // 编辑器宽高是否可变，false-否,true-高可变，'both'-宽高均可，注意引号
  statusbar: true, // 最下方的元素路径和字数统计那一栏是否显示
  elementpath: false, // 元素路径是否显示

  content_style: 'img {max-width:100%;}', // 直接自定义可编辑区域的css样式
  content_css: '/tinymce/skins/content/default/content.css', // 以css文件方式自定义可编辑区域的css样式，css文件需自己创建并引入
  init_instance_callback: (editor) => {

  },
  // images_upload_url: '/demo/upimg.php',  //后端处理程序的url
  // images_upload_base_path: '/demo',  //相对基本路径--关于图片上传建议查看--http://tinymce.ax-z.cn/general/upload-images.php
  // 此处为图片上传处理函数，这个直接用了base64的图片形式上传图片，
  // 如需ajax上传可参考https://www.tiny.cloud/docs/configure/file-image-upload/#images_upload_handler
  images_upload_handler: (blobInfo, success, failure) => {
    const img = 'data:image/jpeg;base64,' + blobInfo.base64()
    success(img)
  },
})

onMounted(() => {
  tinymce.init({})
})

const content = computed({
  get() {
    return props.modelValue
  },
  set(val) {
    context.emit('update:modelValue', val)
  }
})

</script>

<style scoped>
.tinymce-box{
  margin: 20px;
}
</style>
