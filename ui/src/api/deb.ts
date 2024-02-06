 

import { Package, RepositoryType } from './repository'

export const fromDEB = (deb: DEBPackage): Package => ({
  type: RepositoryType.DEB,
  name: deb.name,
  architecture: deb.architecture,
  size: deb.size,
  version: deb.version,
  description: deb.metadata.description,
  projectURL: deb.metadata.projectURL,
  filePath: deb.filePath,
})

export interface DEBPackage {
  name: string
  version: string
  size: number
  architecture: string
  control: string
  metadata: Metadata
  component: string
  distribution: string
  md5: string
  sha1: string
  sha256: string
  sha512: string
  filePath: string
}

export interface Metadata {
  maintainer: string
  projectURL: string
  description: string
  dependencies: string[]
}
