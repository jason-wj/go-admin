import service from '@/utils/request'

// 查询LoginLog列表
export function listLoginLog(query) {
  return service({
    url: '/v1/sys-login-log',
    method: 'get',
    params: query
  })
}

// 查询LoginLog详细
export function getLoginLog(id) {
  return service({
    url: '/v1/sys-login-log/' + id,
    method: 'get'
  })
}

// 删除LoginLog
export function delLoginLog(data) {
  return service({
    url: '/v1/sys-login-log',
    method: 'delete',
    data: data
  })
}

// 导出LoginLog列表
export function exportLoginLog(query) {
  return service({
    url: '/v1/sys-login-log/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

