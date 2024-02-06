 

import { RepositoryType } from '../api/repository'
import { AlpineIcon } from './AlpineIcon'
import { DebianIcon } from './DebianIcon'
import { HelmIcon } from './HelmIcon'
import { RedHatIcon } from './RedHatIcon'

export const packageTypeIcon = (type: RepositoryType) => {
  switch (type) {
    case RepositoryType.DEB:
      return <DebianIcon />
    case RepositoryType.APK:
      return <AlpineIcon />
    case RepositoryType.RPM:
      return <RedHatIcon />
    case RepositoryType.HELM:
      return <HelmIcon />
  }
}
