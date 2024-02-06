 


import { Package, RepositoryType } from './repository'

export const fromHelm = (chart: HelmPackage): Package => ({
  type: RepositoryType.HELM,
  name: chart.name,
  architecture: "noarch",
  size: chart.size,
  version: chart.version,
  description: chart.description,
  projectURL: chart.home || '',
  filePath: chart.filePath,
})

export interface HelmPackage {
  $type: 'helm'
  name: string
  version: string
  description: string
  home: string
  size: number
  filePath: string
}
