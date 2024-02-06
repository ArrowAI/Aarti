 

import { APKPackage, fromAPK } from './apk'
import { DEBPackage, fromDEB } from './deb'
import { fromHelm, HelmPackage } from './helm'
import { fromRPM, RPMPackage } from './rpm'

export interface Stats {
  count: number
  size: number
}

export interface Repository {
  name?: string
  type: RepositoryType
  size: number
  lastUpdated: Date
  metadata: Stats
  packages: Stats
}

export enum RepositoryType {
  APK = 'apk',
  DEB = 'deb',
  RPM = 'rpm',
  HELM = 'helm',
}

export interface Package {
  type: RepositoryType
  name: string
  size: number
  version: string
  architecture: string
  license?: string
  description: string
  summary?: string
  projectURL: string
  lastUpdated?: Date
  filePath: string
}

export const makePackage = (type: RepositoryType, d: APKPackage | DEBPackage | RPMPackage | HelmPackage) => {
  switch (type) {
    case RepositoryType.APK:
      return fromAPK(d as APKPackage)
    case RepositoryType.DEB:
      return fromDEB(d as DEBPackage)
    case RepositoryType.RPM:
      return fromRPM(d as RPMPackage)
    case RepositoryType.HELM:
      return fromHelm(d as HelmPackage)
  }
}

export const subRepositories = (packages: Package[], type: RepositoryType) => type !== RepositoryType.RPM && type !== RepositoryType.HELM
  ? packages.map(p => p.filePath.replace('pool/', '').split('/').slice(0, 2).join('/')).filter((p, i, arr) => arr.indexOf(p) === i)
  : []

export const subRepositoryPackages = (packages: Package[], type: RepositoryType, sub: string) => packages.filter(p => sub === '' || (type === RepositoryType.DEB ? p.filePath.startsWith(`pool/${sub}`) : p.filePath.startsWith(sub)))
