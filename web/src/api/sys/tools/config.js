import service from '@/utils/request'

// @Tags system
// @Summary è·ċéuĉ
export const getConfig = () => {
  return service({
    url: '/v1/sysRuntimeConfig/getConfig',
    method: 'get',
    donNotShowLoading: true
  })
}
