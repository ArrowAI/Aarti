 


import { Package, RepositoryType } from './repository'

export const fromRPM = (rpm: RPMPackage): Package => ({
  type: RepositoryType.RPM,
  name: rpm.name,
  architecture: rpm.fileMetadata.architecture,
  lastUpdated: new Date(rpm.fileMetadata.buildTime * 1000),
  size: rpm.size,
  version: rpm.version,
  license: rpm.versionMetadata.license || '',
  description: rpm.versionMetadata.description,
  summary: rpm.versionMetadata.summary,
  projectURL: rpm.versionMetadata.projectURL || '',
  filePath: rpm.filePath,
})

export interface RPMPackage {
  $type: 'rpm'
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
