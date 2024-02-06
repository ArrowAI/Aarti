 

import { CheckBoxOutlineBlank } from '@mui/icons-material'
import React from 'react'
import { RepositoryType } from '../api/repository'
import PackagesPage from '../Components/Packages/PackagesPage'
import { MainRoutesRegistry } from '../routes'

MainRoutesRegistry['packages'] = {
  path: '/:repo',
  component: <PackagesPage />,
  icon: <CheckBoxOutlineBlank />,
  priority: 200,
  public: false,
  label: 'Packages',
  show: false,
  navigate: ([type, repo]: [type: RepositoryType, repo: string]) => `/${type}/${repo}`,
}
