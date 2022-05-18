import service from '@/utils/request'

// 查询服务器详细
export function getServer() {
  return service({
    url: '/v1/server-monitor',
    method: 'get'
  })
}
