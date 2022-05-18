import { formatTimeToStr } from '@/utils/date'
import i18n from '@/i18n' // added by mohamed hassan to multilangauge

export const formatBoolean = (bool) => {
  if (bool !== null) {
    return bool ? i18n.t('general.yes') : i18n.t('general.no')
  } else {
    return ''
  }
}
export const formatDate = (time) => {
  if (time !== null && time !== '') {
    var date = new Date(time)
    return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
  } else {
    return ''
  }
}

