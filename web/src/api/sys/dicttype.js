import service from '@/utils/request'

// 查询字典类型列表
export function listSysDictType(query) {
  return service({
    url: '/v1/dict/type',
    method: 'get',
    params: query
  })
}

// 查询字典类型详细
export function getSysDictType(dictId) {
  return service({
    url: '/v1/dict/type/' + dictId,
    method: 'get'
  })
}

// 新增字典类型
export function addSysDictType(data) {
  return service({
    url: '/v1/dict/type',
    method: 'post',
    data: data
  })
}

// 修改字典类型
export function updateSysDictType(data) {
  return service({
    url: '/v1/dict/type/' + data.dictId,
    method: 'put',
    data: data
  })
}

// 删除字典类型
export function delSysDictType(dictId) {
  return service({
    url: '/v1/dict/type',
    method: 'delete',
    data: dictId
  })
}

// 导出字典类型
export function exportSysDictType(query) {
  return service({
    url: '/v1/dict/type/export',
    method: 'get',
    responseType: 'blob',
    params: query
  })
}

// 获取字典选择框列表
export function optionselect() {
  return service({
    url: '/v1/dict/type-option-select',
    method: 'get'
  })
}
