import service from '@/utils/request'

// 查询SysOperlog列表
export function listOperaLog(query) {
  return service({
    url: '/v1/sys-opera-log',
    method: 'get',
    params: query
  })
}

// 删除SysOperlog
export function delOperalog(data) {
  return service({
    url: '/v1/sys-opera-log',
    method: 'delete',
    data: data
  })
}

// 导出OperaLog列表
export function exportOperaLog(query) {
  return service({
    url: '/v1/sys-opera-log/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

