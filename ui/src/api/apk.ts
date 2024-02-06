 


import { Package, RepositoryType } from './repository'

export const fromAPK = (apk: APKPackage): Package => ({
  type: RepositoryType.APK,
  name: apk.name,
  architecture: apk.fileMetadata.architecture,
  size: apk.size,
  version: apk.version,
  license: apk.versionMetadata.license || '',
  description: apk.versionMetadata.description,
  summary: apk.versionMetadata.summary,
  projectURL: apk.versionMetadata.projectURL || '',
  filePath: apk.filePath,
})

export interface APKPackage {
  $type: 'apk'
  name: string
  version: string
  versionMetadata: VersionMetadata
  fileMetadata: FileMetadata
  hashSha256: string
  size: number
  filePath: string
}

export interface VersionMetadata {
  license?: string
  projectURL?: string
  summary: string
  description: string
}

export interface FileMetadata {
  architecture: string
  epoch: string
  version: string
  release: string
  vendor?: string
  group?: string
  packager: string
  sourceRPM: string
  buildHost: string
  buildTime: number
  fileTime?: number
  installedSize?: number
  archiveSize: number
  provide: Provide[]
  require?: Require[]
  files?: File[]
  changelogs?: Changelog[]
  conflict?: Conflict[]
  obsolete?: Obsolete[]
}

export interface Provide {
  name: string
  flags?: string
  version?: string
  epoch?: string
  release?: string
}

export interface Require {
  name: string
  flags?: string
  version?: string
  epoch?: string
  release?: string
}

export interface File {
  path: string
  isExecutable: boolean
  type?: string
}

export interface Changelog {
  author: string
  date: number
  text: string
}

export interface Conflict {
  name: string
  flags: string
  version: string
  epoch: string
  release?: string
}

export interface Obsolete {
  name: string
  flags: string
  version: string
  epoch: string
  release?: string
}
