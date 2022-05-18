import service from '@/utils/request'

// 查询Config列表
export function listConfig(query) {
  return service({
    url: '/v1/config',
    method: 'get',
    params: query
  })
}

// 查询Config详细
export function getConfig(id) {
  return service({
    url: '/v1/config/' + id,
    method: 'get'
  })
}

// 新增Config
export function addConfig(data) {
  return service({
    url: '/v1/config',
    method: 'post',
    data: data
  })
}

// 修改Config
export function updateConfig(data) {
  return service({
    url: '/v1/config/' + data.id,
    method: 'put',
    data: data
  })
}

// 删除Config
export function delConfig(data) {
  return service({
    url: '/v1/config',
    method: 'delete',
    data: data
  })
}

// 导出Config列表
export function exportConfig(query) {
  return service({
    url: '/v1/config/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

// 根据参数键名查询参数值
export function getConfigKey(configKey) {
  return service({
    url: '/v1/configKey/' + configKey,
    method: 'get'
  })
}

