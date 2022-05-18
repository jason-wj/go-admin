/**
 * 网站配置文件
 */

const config = {
  appName: 'Go-Admin-Web',
  appLogo: 'http://www.wjblog.top/images/my_head-touch-icon-next.png',
  showViteLogo: true
}

export const viteLogo = (env) => {
  if (config.showViteLogo) {
    const chalk = require('chalk')
    console.log(
      chalk.green(
        `> 欢迎使用Gin-Vue-Admin，开源地址：https://github.com/flipped-aurora/gin-vue-admin`
      )
    )
    console.log(
      chalk.green(
        `> 当前版本:V2.5.0`
      )
    )
    console.log('\n')
  }
}

export default config
