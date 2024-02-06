 

import { PowerSettingsNewOutlined } from '@mui/icons-material'
import { useNavigate } from 'react-router-dom'
import { useAPI } from '../api/useAPI'
import { useAsyncOnce } from '../hooks'
import { MainRoutesRegistry } from '../routes'

export const LogoutPage = () => {
  const navigate = useNavigate()
  const { logout, authenticated } = useAPI()
  useAsyncOnce(async () => {
    if (authenticated) {
      await logout()
    }
    navigate(MainRoutesRegistry.login.navigate(), { replace: true })
  })
  return null
}

MainRoutesRegistry['logout'] = {
  label: 'Logout',
  path: '/logout',
  component: <LogoutPage />,
  icon: <PowerSettingsNewOutlined />,
  priority: 0,
  public: false,
  show: false,
  bottomEnd: true,
  navigate: () => '/logout',
}
