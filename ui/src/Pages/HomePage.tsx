 

import { HomeOutlined } from '@mui/icons-material'
import React from 'react'
import HomePage from '../Components/Home/HomePage'
import { MainRoutesRegistry } from '../routes'

MainRoutesRegistry['home'] = {
  path: '/',
  component: <HomePage />,
  icon: <HomeOutlined />,
  priority: 100,
  public: false,
  label: 'Home',
  show: true,
  navigate: () => '/',
}
