import service from '@/utils/request'

export function getDeptList(query) {
  return service({
    url: '/v1/sys/dept',
    method: 'get',
    params: query
  })
}

// 查询部门详细
export function getDept(deptId) {
  return service({
    url: '/v1/sys/dept/' + deptId,
    method: 'get'
  })
}

// 查询部门下拉树结构
export function treeselect() {
  return service({
    url: '/v1/sys/deptTree',
    method: 'get'
  })
}

// 根据角色ID查询部门树结构
export function roleDeptTreeselect(roleId) {
  return service({
    url: '/v1/sys/roleDeptTreeselect/' + roleId,
    method: 'get'
  })
}

// 新增部门
export function addDept(data) {
  return service({
    url: '/v1/sys/dept',
    method: 'post',
    data: data
  })
}

// 修改部门
export function updateDept(data, id) {
  return service({
    url: '/v1/sys/dept/' + id,
    method: 'put',
    data: data
  })
}

// 删除部门
export function delDept(deptId) {
  return service({
    url: '/v1/sys/dept/' + deptId,
    method: 'delete'
  })
}
