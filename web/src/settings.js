export default {
  title: 'Go-Admin后台管理系统',

  /**
     * @type {boolean} true | false
     * @description Whether show the settings right-panel
     */
  showSettings: true,

  /**
     * 是否显示顶部导航
     */
  topNav: false,

  /**
     * @type {boolean} true | false
     * @description Whether need tagsView
     */
  tagsView: false,

  /**
     * @type {boolean} true | false
     * @description Whether fix the header
     */
  fixedHeader: true,

  /**
     * @type {boolean} true | false
     * @description Whether show the logo in sidebar
     */
  sidebarLogo: true,

  /**
     * @type {string | array} 'production' | ['production', 'development']
     * @description Need show err logs component.
     * The default is only used in the production env
     * If you want to also use it in dev, you can pass ['production', 'development']
     */
  errorLog: 'production',

  /**
     * dark - 黑色 light-浅色   #fff - 自定义颜色
     * */
  sideMode: 'dark',

  /**
     * 活跃颜色
     * */
  activeColor: '#4D70FF',

  /**
     * 默认基础颜色
     * */
  baseColor: '#fff',

  /**
     * 背景色
     * */
  backgroundColor: '#191a23',

  /**
     * 默认路由
     * */
  defaultRouter: 'dashboard'
}
