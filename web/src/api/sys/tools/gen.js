import service from '@/utils/request'

// 查询生成表数据
export function listTable(query) {
  return service({
    url: '/v1/sys/tables/page',
    method: 'get',
    params: query
  })
}
// 查询db数据库列表
export function listDbTable(query) {
  return service({
    url: '/v1/db/tables/page',
    method: 'get',
    params: query
  })
}

// 查询表详细信息
export function getGenTable(tableId) {
  return service({
    url: '/v1/sys/tables/info/' + tableId,
    method: 'get'
  })
}

export function getGenTableInfo(tablename) {
  return service({
    url: '/v1/sys/tables?tableName=' + tablename,
    method: 'get'
  })
}

// 修改代码生成信息
export function updateGenTable(data) {
  return service({
    url: '/v1/sys/tables/info',
    method: 'put',
    data: data
  })
}

// 导入表
export function importTable(data) {
  return service({
    url: '/v1/sys/tables/info',
    method: 'post',
    params: data
  })
}
// 预览生成代码
export function previewTable(tableId) {
  return service({
    url: '/v1/gen/preview/' + tableId,
    method: 'get'
  })
}
// 删除表数据
export function delTable(tableId) {
  return service({
    url: '/v1/sys/tables/info/' + tableId,
    method: 'delete'
  })
}

// 生成代码到项目
export function toProjectTable(tableId) {
  return service({
    url: '/v1/gen/toproject/' + tableId,
    method: 'get'
  })
}

// 下载代码
export function downloadCode(tableId) {
  return service({
    url: '/v1/gen/downloadCode/' + tableId,
    responseType: 'blob',
    method: 'get'
  })
}

// 生成菜单到数据库
export function toDBTable(tableId) {
  return service({
    url: '/v1/gen/todb/' + tableId,
    method: 'get'
  })
}

export function getTableTree() {
  return service({
    url: '/v1/gen/tabletree',
    method: 'get'
  })
}
